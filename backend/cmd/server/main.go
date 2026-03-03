package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"FGDSchemaBuilder/internal/server"
)

func main() {
	// Crear el multiplexer (router)
	mux := http.NewServeMux()

	// Registrar rutas de forma explícita pasándole el mux
	server.RegisterRoutes(mux)

	slog.Log(context.TODO(), slog.LevelInfo, "Server started", "port", "8000")

	// Escuchar y servir usando el mux local
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
