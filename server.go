package main

import (
	"encoding/json"
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

// WebSocketRequest represents a request made from the client via WebSockets. It includes
// the type of request so that we know what handler function to use with it and
// the data associated with the request, which is passed on to  the handler function.
// At the moment, all keys and values in the request data are strings.
type WebSocketRequest struct {
	RequestType string            `json:"type"`
	RequestData map[string]string `json:"data"`
}

// NewServer attempts to open the given database file and returns a new Server if
// successful.
func NewServer(dbPath string) (*Server, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
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

	if r.URL.Path == "/debug" {
		http.ServeFile(w, r, "debug.html")
		return
	}

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
	defer ws.Close()

	// TODO: insert infinite for loop that ends when the connection is closed for reading and handling requests
	keepListening := true

	for keepListening {
		_, requestText, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading client request from websocket: %v\n", err)
			break
		}

		var requestObject WebSocketRequest
		if err = json.Unmarshal(requestText, &requestObject); err != nil {
			log.Printf("Invalid request JSON received: %v", err)
			return
		}

		switch requestObject.RequestType {
		case "login":
			{
				loginSuccess := HandleLogin(requestObject, s)
				if !loginSuccess {
					return
				}
			}
		default:
			{
				fmt.Printf("Error finding request type %v", requestObject.RequestType)
			}
		}
	}

	//Remnants of old code that haven't been cleaned up yet
	// var user UserInfo
	// s.Database.Take(&user, &UserInfo{Username: loginInfo.Username})
}
