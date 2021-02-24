package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Server represents this server as a whole and contains global configuration
// information so request-handling code has a single spot to read it from.
type Server struct {
	// The database where all our persistent information will be stored, i.e.,
	// basically everything. Under the hood, the information will get stored as
	// a file on the disk because we will be using SQLite for now. That may
	// change in the future.
	Database *gorm.DB
}

// NewServer attempts to open the given database file and returns a new Server if
// successful.
func NewServer(dbPath string) (*Server, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	return &Server{db}, nil
}

// ListenAndServe begins listening on the given local address, accepting only WebSocket
// connections at '/connect'. This method will block until an error occurs and the server
// is no longer listening and serving HTTP requests at the address.
func (s *Server) ListenAndServe(addr string) error {
	log.Printf("Listening on '%s'...", addr)
	return http.ListenAndServe(addr, s)
}

// ServeHTTP implements the http.Handler interface for Server. The only HTTP route provided
// is '/connect', which immediately upgrades request connections to WebSockets. A proprietary
// WebSocket protocol is used from there.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We only care about /connect requests. Returning early saves us from annoying code indentation
	// below.
	if r.URL.Path != "/connect" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// The .Upgrade() call promises to send an error reply to the browser if upgrading fails so
		// we can just log the error here and return.
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	// INSERT AUTHENTICATION CODE HERE

	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error reading message from websocket: %v\n", err)
			break
		}

		fmt.Printf("message received: %d %s\n", messageType, string(data))
	}
}
