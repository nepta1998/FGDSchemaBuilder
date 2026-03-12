// Package service contains the business logic.
package service

import (
	"FGDSchemaBuilder/internal/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// ParserService maneja la lógica de negocio para procesar archivos FGD
type ParserService struct{}

// NewParserService crea una nueva instancia del servicio
func NewParserService() *ParserService {
	return &ParserService{}
}

func (s *ParserService) ParseFGD(fgdText string) models.FGD {
	clenedText := s.cleanFGD(fgdText)
	mapSize, _ := s.parseMapSize(clenedText)
	includes := s.parseIncludes(clenedText)

	metadata := models.Metadata{
		MapSize:  mapSize,
		Includes: includes,
	}
	entities := s.parseEntityClasses(clenedText)

	fgd := models.FGD{
		Metadata: metadata,
		Entities: entities,
	}

	return fgd
}

func (s *ParserService) cleanFGD(fgdText string) string {
	re := regexp.MustCompile(`(?s:/\*.*?\*/)|//.*`)
	cleanText := re.ReplaceAllString(fgdText, "")
	cleanText = strings.ReplaceAll(cleanText, "\r\n", "\n")
	cleanText = strings.TrimSpace(cleanText)
	return cleanText
}

func (s *ParserService) parseMapSize(cleanedText string) (*models.MapSize, error) {
	re := regexp.MustCompile(`@mapsize\s*\(\s*(-?\d+)\s*,\s*(-?\d+)\s*\)`)
	match := re.FindStringSubmatch(cleanedText)
	if match != nil {
		min, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("invalid min value: %v", err)
		}
		max, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, fmt.Errorf("invalid max value: %v", err)
		}
		return &models.MapSize{
			Min: min,
			Max: max,
		}, nil
	} else {
		return nil, fmt.Errorf("could not find @mapsize")
	}
}

func (s *ParserService) parseIncludes(cleanedText string) []string {
	re := regexp.MustCompile(`@include\s*"([^"]+)"`)
	matches := re.FindAllStringSubmatch(cleanedText, -1)

	var includes []string
	for _, match := range matches {
		if len(match) > 1 {
			includes = append(includes, match[1])
		}
	}
	return includes
}

func (s *ParserService) parseEntityClasses(cleanedText string) []models.Entity {
	const pattern = `@(\w+)\s*([\s\S]*?)\s*=\s*(\w+)\s*(?::\s*"([^"]*)")?\s*(?:\[([\s\S]*?)\])?`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(cleanedText, -1)

	entities := make([]models.Entity, 0, len(matches))
	for _, m := range matches {
		if len(m) < 6 {
			continue
		}
		classType := m[1]
		header := m[2]
		name := m[3]
		description := m[4]
		body := m[5]

		entity := models.Entity{
			ID:          uuid.New().String(),
			ClassType:   classType,
			Name:        name,
			Description: description,
			BaseClasses: s.parseBaseClasses(header),
			Helpers:     s.parseHelpers(header),
			Properties:  s.parseProperties(body),
		}
		entities = append(entities, entity)
	}
	return entities
}

func (s *ParserService) parseBaseClasses(baseClassesText string) []string {
	if baseClassesText == "" {
		return []string{}
	}
	const pattern = `base\(([^)]+)\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(baseClassesText, -1)

	var bases []string
	for _, match := range matches {
		if len(match) > 1 {
			for className := range strings.SplitSeq(match[1], ",") {
				trimmed := strings.TrimSpace(className)
				if trimmed != "" {
					bases = append(bases, trimmed)
				}
			}
		}
	}
	return bases
}

func (s *ParserService) parseHelpers(helpersText string) models.Helpers {
	re := regexp.MustCompile(`(size|color)\s*\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(helpersText, -1)
	helpers := models.Helpers{}
	for _, match := range matches {
		helpers[match[1]] = strings.TrimSpace(match[2])
	}
	return helpers
}

func (s *ParserService) parseProperties(propertiesText string) []models.Property {
	lines := strings.Split(propertiesText, "\n")
	cleanLines := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			cleanLines = append(cleanLines, trimmed)
		}
	}
	properties := make([]models.Property, 0, len(cleanLines))
	propRegex := regexp.MustCompile(`(\w+)\s*\((\w+)\)\s*(?::\s*"([^"]*)")?\s*(?::\s*(-?[\d.\s]+|"[^"]*"))?\s*(?::\s*"([^"]*)")?`)
	blockContent := ""
	for i := 0; i < len(cleanLines); i++ {
		line := cleanLines[i]
		isBlockProp := strings.Contains(line, "=") &&
			(strings.Contains(line, "(flags)") || strings.Contains(line, "(choices)"))
		if isBlockProp {
			for j := i + 1; j < len(cleanLines); j++ {
				before, _, wasFind := strings.Cut(cleanLines[j], "]")
				if wasFind {
					blockContent += before
					i = j
					break
				}
				blockContent += cleanLines[j] + "\n"
			}
		}
		if match := propRegex.FindStringSubmatch(line); match != nil {
			name := match[1]
			typeMatch := match[2]
			displayName := match[3]
			defaultValue := match[4]
			description := match[5]
			options := []models.Option{}
			if typeMatch == "" {
				typeMatch = "string"
			}
			if blockContent != "" {
				isFlags := strings.ToLower(typeMatch) == "flags"
				if isFlags {
					options = s.parseFlags(fmt.Sprintf("[%s]", blockContent))
				} else {
					options = s.parseChoices(blockContent)
				}
				blockContent = ""
			}
			if defaultValue != "" {
				defaultValue = strings.TrimSpace(strings.ReplaceAll(defaultValue, "\"", ""))
			}
			property := models.Property{
				ID:           uuid.New().String(),
				Name:         name,
				Type:         typeMatch,
				DisplayName:  displayName,
				DefaultValue: defaultValue,
				Description:  description,
				Options:      options,
			}
			properties = append(properties, property)
		}
	}
	return properties
}

func (s *ParserService) parseFlags(flagsLine string) []models.Option {
	re := regexp.MustCompile(`(\w+)\s*=\s*(-?\d+)`)
	match := re.FindStringSubmatch(flagsLine)
	if match == nil {
		return []models.Option{}
	}
	re = regexp.MustCompile(`(\d+)\s*:\s*"([^"]+)"\s*:\s*(\d)`)
	matches := re.FindAllStringSubmatch(flagsLine, -1)
	options := make([]models.Option, 0, len(matches))
	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		label := match[2]
		isDefault := strings.TrimSpace(match[3]) == "1"
		option := models.Option{
			Value:   value,
			Label:   label,
			Default: isDefault,
		}
		options = append(options, option)
	}
	return options
}

func (s *ParserService) parseChoices(choiceBlock string) []models.Option {
	re := regexp.MustCompile(`(-?\d+)\s*:\s*"([^"]+)"`)
	matches := re.FindAllStringSubmatch(choiceBlock, -1)
	options := make([]models.Option, 0, len(matches))
	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		label := match[2]

		option := models.Option{
			Value: value,
			Label: label,
		}
		options = append(options, option)
	}
	return options
}
