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

func ParseFGD(fgdText string) models.FGD {
	clenedText := cleanFGD(fgdText)
	mapSize, _ := parseMapSize(clenedText)
	includes := parseIncludes(clenedText)

	metadata := models.Metadata{
		MapSize:  mapSize,
		Includes: includes,
	}
	entities := parseEntityClasses(clenedText)

	fgd := models.FGD{
		Metadata: metadata,
		Entities: entities,
	}

	return fgd
}

func cleanFGD(fgdText string) string {
	// 1. Remove comments (multiline /* */ and single line //)
	// In Go, (?s) is the flag to make the dot (.) match newlines
	re := regexp.MustCompile(`(?s)/\*.*?\*/|//.*`)
	cleanText := re.ReplaceAllString(fgdText, "")
	// 2. Normalize line endings (Windows \r\n -> Unix \n)
	cleanText = strings.ReplaceAll(cleanText, "\r\n", "\n")
	// 3. Trim leading and trailing whitespace
	cleanText = strings.TrimSpace(cleanText)
	return cleanText
}

func parseMapSize(cleanedText string) (*models.MapSize, error) {
	// Definimos el regex. Usamos backticks `` para que sea un raw string
	// y no tener que escapar las barras invertidas.
	re := regexp.MustCompile(`@mapsize\s*\(\s*(-?\d+)\s*,\s*(-?\d+)\s*\)`)
	// FindStringSubmatch devuelve un slice:
	// match[0] es la coincidencia completa
	// match[1] es el primer grupo (-?\d+)
	// match[2] es el segundo grupo (-?\d+)
	match := re.FindStringSubmatch(cleanedText)
	if match != nil {
		min, err := strconv.Atoi(match[1])
		if err != nil {
			return &models.MapSize{}, fmt.Errorf("invalid min value: %v", err)
		}
		max, err := strconv.Atoi(match[2])
		if err != nil {
			return &models.MapSize{}, fmt.Errorf("invalid max value: %v", err)
		}
		return &models.MapSize{
			Min: min,
			Max: max,
		}, nil
	} else {
		return &models.MapSize{}, fmt.Errorf("could not find @mapsize")
	}
}

func parseIncludes(cleanedText string) []string {
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

func parseEntityClasses(cleanedText string) []models.Entity {
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
			BaseClasses: parseBaseClasses(header),
			Helpers:     parseHelpers(header),
			Properties:  parseProperties(body),
		}
		entities = append(entities, entity)
	}
	return entities
}

func parseBaseClasses(baseClassesText string) []string {
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

func parseHelpers(helpersText string) models.Helpers {
	re := regexp.MustCompile(`(size|color)\s*\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(helpersText, -1)
	helpers := models.Helpers{}
	for _, match := range matches {
		helpers[match[1]] = strings.TrimSpace(match[2])
	}
	return helpers
}

func parseProperties(propertiesText string) []models.Property {
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
				before, _, wasFind := strings.Cut(line, "]")
				if wasFind {
					blockContent += before
					i = j
					break
				}
				blockContent += line + "\n"
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
					options = parseFlags(fmt.Sprintf("[%s]", blockContent))
				} else {
					options = parseChoices(blockContent)
					blockContent = ""
				}
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

func parseFlags(flagsLine string) []models.Option {
	return []models.Option{}
}

func parseChoices(choiceBlock string) []models.Option {
	return []models.Option{}
}
