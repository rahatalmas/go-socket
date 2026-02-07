package main

import (
	"butter-socket/internal/handler"
	"butter-socket/internal/hub"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting WebSocket Server...")

	// Create and start the hub
	h := hub.NewHub()
	go h.Run()

	// Setup routes
	http.HandleFunc("/ws/customer", func(w http.ResponseWriter, r *http.Request) {
		handler.WsHandler(h, w, r)
	})

	http.HandleFunc("/ws/user", func(w http.ResponseWriter, r *http.Request) {
		handler.WsUserHandler(h, w, r)
	})

	// Start server
	addr := ":4646"
	fmt.Printf("Server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
