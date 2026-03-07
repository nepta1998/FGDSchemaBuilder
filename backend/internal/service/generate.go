package service

import (
	"fmt"
	"strings"

	"FGDSchemaBuilder/internal/models"
)

type GenerateService struct{}

// NewGenerateService crea una nueva instancia del servicio
func NewGenerateService() *GenerateService {
	return &GenerateService{}
}

func (s *GenerateService) GenerateFGD(schema *models.FGD) string {
	output := []string{}
	if schema.Metadata.MapSize != nil {
		output = append(
			output,
			fmt.Sprintf(
				"@mapsize(%d, %d)",
				schema.Metadata.MapSize.Min,
				schema.Metadata.MapSize.Max,
			),
		)
	}
	if len(schema.Metadata.Includes) > 0 {
		for _, include := range schema.Metadata.Includes {
			output = append(output, fmt.Sprintf(`@include "%s"`, include))
		}
	}

	return strings.Join(output, "\n\n")
}
