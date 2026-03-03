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

	// 2. Inicializar los manejadores e inyectar sus servicios
	parserHdl := handler.NewParserHandler(parserSvc)

	// 3. Registrar las rutas
	mux.HandleFunc("POST "+PathParse, parserHdl.Parse)
}
