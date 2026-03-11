package server

import (
	"FGDSchemaBuilder/internal/handler"
	"FGDSchemaBuilder/internal/service"
	"net/http"
)

const (
	PathParse    = "/parse"
	PathGenerate = "/generate"
)

func RegisterRoutes(mux *http.ServeMux) {
	// 1. Inicializar los servicios (las herramientas)
	parserSvc := service.NewParserService()
	generateSvc := service.NewGenerateService()

	// 2. Inicializar los manejadores e inyectar sus servicios
	parserHdl := handler.NewParserHandler(parserSvc)
	generateHdl := handler.NewGenerateHandler(generateSvc)


	// 3. Registrar las rutas
	mux.HandleFunc("POST "+PathParse, parserHdl.Parse)
	mux.HandleFunc("POST "+PathGenerate, generateHdl.Generate)
}
