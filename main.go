package main

import (
	"log"
	"os"
)

func main() {
	server, err := NewServer("debug/cache.sqlite3")
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err = server.ListenAndServe(":8080"); err != nil {
		log.Fatalf("Server quit listening unexpectedly: %v", err)
	}

	os.Exit(0)
}
