// Package service contains the business logic.
package service

import (
	"FGDSchemaBuilder/internal/models"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParseFGD(fgdText string) (models.FGD, error) {
	clenedText := cleanFGD(fgdText)
	mapSize, _ := getMapSize(clenedText)
	includes := getInclude(clenedText)

	metadata := models.Metadata{
		MapSize:  mapSize,
		Includes: includes,
	}

	fgd := models.FGD{
		Metadata: metadata,
	}

	return fgd, nil
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

func getMapSize(cleanedText string) (*models.MapSize, error) {
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

func getInclude(cleanedText string) []string {
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
