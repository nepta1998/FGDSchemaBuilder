package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"FGDSchemaBuilder/internal/server"
	"FGDSchemaBuilder/internal/server/mainroute"
)

func InitRoutes() {
	mainroute.ParseView()
}

func main() {
	InitRoutes()
	slog.Log(context.TODO(), slog.LevelInfo, "Server started", "port", "8000")
	err := http.ListenAndServe(":8000", server.GetMuxInstance())
	if err != nil {
		log.Fatal(err)
	}
}

// func main() {
// 	content, err := os.ReadFile("/home/neptali/projects/FGDSchemaBuilder/Quake.fgd")
// 	if err != nil {
// 		fmt.Println("Error al abrir el archivo:", err)
// 		return
// 	}
//
// 	text := string(content)
// 	fgdParse := service.ParseFGD(text)
// 	jsonData, err := json.MarshalIndent(fgdParse, "", "  ")
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	err = os.WriteFile("output.json", jsonData, 0644)
// 	if err != nil {
// 		fmt.Println("Error al guardar el archivo:", err)
// 		return
// 	}
//
// 	fmt.Println("✅ JSON guardado exitosamente en output.json")
// }
