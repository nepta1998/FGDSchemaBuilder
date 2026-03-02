package mainroute

import (
	"FGDSchemaBuilder/internal/handler"
	"FGDSchemaBuilder/internal/server"
)

func ParseView() {
	server.GetMuxInstance().HandleFunc(server.RoutesInstance.PARSE, handler.ParseHandler)
}
