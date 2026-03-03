package handler

import (
	"FGDSchemaBuilder/internal/service"
	"encoding/json"
	"io"
	"net/http"
)

// ParserHandler agrupa todos los manejadores relacionados con el parseo de archivos FGD
type ParserHandler struct {
	parserService *service.ParserService
}

// NewParserHandler es el constructor para ParserHandler, recibe el servicio como dependencia
func NewParserHandler(svc *service.ParserService) *ParserHandler {
	return &ParserHandler{
		parserService: svc,
	}
}

// Parse maneja la petición POST para procesar el contenido de un archivo FGD
func (h *ParserHandler) Parse(w http.ResponseWriter, r *http.Request) {
	// 1. Leer el cuerpo de la petición (el texto del archivo FGD)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la petición", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 2. Llamar al servicio inyectado
	fgdResult := h.parserService.ParseFGD(string(body))

	// 3. Responder con el JSON resultante
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fgdResult)
}
