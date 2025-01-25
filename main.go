package main

import (
	"fmt"
	"log"
	"net/http"
	"websocket-server/database"
	"websocket-server/routes"
)

func main() {

	database.InitializePostgresDB()

	mux := http.NewServeMux()

	// Register routes
	// routes.RegisterAuthRoutes(mux)

	routes.RegisterUserRoutes(mux)
	routes.RegisterMessagingRoutes(mux)

	// WebSocket endpoint
	// http.HandleFunc("/ws", handlers.WebSocketHandler)

	// Token generation endpoint
	// http.HandleFunc("/generate-token", handlers.GenerateTokenHandler)

	port := "10000"
	fmt.Printf("Server started on port %s\n", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Fatal(http.ListenAndServe(":"+port, mux))

}
