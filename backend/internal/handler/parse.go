package handler

import (
	"net/http"

	"FGDSchemaBuilder/internal/server"
	"FGDSchemaBuilder/internal/utils"
)

func ParseHandler(w http.ResponseWriter, r *http.Request) {
	utils.CheckIfPath(w, r, server.RoutesInstance.PARSE)
}
