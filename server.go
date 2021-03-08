package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
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
	ID   uint   `json:"id"`   // request ID (counter) for mapping back a response
	Type string `json:"type"` // what type of request is it?
	Data string `json:"data"` // request body; shape depends on the request type
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

	type LoginRequest struct {
		Username string `json:"username"`
		Password []byte `json:"password"` // TODO: we NEED to use msgpack, JSON will expect Base64 for the password now lol
	}

	// TODO: add a 5s timeout so we don't sit and wait forever for a login request that
	// might not be coming.

	_, loginReqJSON, err := ws.ReadMessage()
	if err != nil {
		log.Printf("Failed to read login request: %v", err)
		return
	}

	var loginReq LoginRequest
	err = json.Unmarshal(loginReqJSON, &loginReq)
	scrub(loginReqJSON) // wipe the password from memory
	if err != nil {
		log.Printf("Login request was not valid JSON: %v", err)
		return
	}
	log.Printf("Login request received for %q", loginReq.Username)

	var user UserInfo
	if err = s.Database.Take(&user, &UserInfo{Username: loginReq.Username}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User does not exist: %q", loginReq.Username)
		} else {
			log.Printf("Failed to lookup %q from database: %v", loginReq.Username, err)
		}
		scrub(loginReq.Password)
		return // TODO: send back some sort of error feedback to the client
	}
	log.Printf("Found %q in database, attempting to authenticate", loginReq.Username)

	saltedPassword := append(user.Salt, loginReq.Password...)
	scrub(loginReq.Password)
	passwordHash := sha256.Sum256(saltedPassword)
	scrub(saltedPassword)

	// Start i at 1 because we already did 1 round of hashing to create the variable above
	for i := uint(1); i < user.Rounds; i++ {
		passwordHash = sha256.Sum256(passwordHash[:])
	}

	if !bytes.Equal(passwordHash[:], user.PasswordHash) {
		log.Printf("Authentication failed for %q: incorrect password", loginReq.Username)
		return
	}
	log.Printf("Authentication successful for %q", loginReq.Username)

	(&UserConn{
		UserInfo: user,
		Conn:     ws,
	}).Serve()
}
