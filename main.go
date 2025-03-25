package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Flags
	port := flag.Int("port", 8080, "Port number for the server (1-65535)")
	storageDir := flag.String("dir", "storage", "Path to the storage directory")
	flag.Parse()

	// Validation
	if *port <= 0 || *port > 65535 {
		log.Fatalf("Invalid port number: %d", *port)
	}
	if *storageDir == "" {
		log.Fatal("Storage directory path is required")
	}
	if _, err := os.Stat(*storageDir); os.IsNotExist(err) {
		log.Fatalf("Storage directory does not exist: %s", *storageDir)
	}

	// Handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("well well well first step is done ahhahaha"))
	})

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting server on %s...\n", addr[1:])
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error 500. Failed to start server: %v", err)
	}
}
