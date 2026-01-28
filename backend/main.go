package main

import (
	"log"
	"net/http"
	"social-network/app/handlers/websocket"
	"social-network/server"

	"social-network/db"
)

func main() {
	// initializing the database:
	err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Database.Close()

	log.Println("Database initialized successfully")

	//start hub for websocket:
	go websocket.StartHub()

	// making the outside-wrapper handler with CORS enabled:
	router := server.SetupRoutes()

	// starting the server:
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
