package mainroute

import (
	"FGDSchemaBuilder/internal/handler"
	"FGDSchemaBuilder/internal/server"
)

func ParseView() {
	server.GetMuxInstance().HandleFunc("POST "+server.RoutesInstance.PARSE, handler.ParseHandler)
}
