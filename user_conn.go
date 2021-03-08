package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// UserConn is an authenticated websocket connection.
type UserConn struct {
	UserInfo
	*websocket.Conn
}

// Serve blocks, repeatedly reading and handling individual requests asynchronously until reading
// a message from the websocket fails.
func (c *UserConn) Serve() error {
	for {
		_, reqJSON, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading client request from websocket: %v\n", err)
			return err
		}

		var req WebSocketRequest
		if err = json.Unmarshal(reqJSON, &req); err != nil {
			log.Printf("Invalid request JSON received: %v", err)
			// TODO: should we close the websocket? How do we map back an error reply when
			// the request was malformed and we don't even have a request ID? Do we just
			// continue ignoring malformed requests and act like nothing happened?
			continue
		}

		switch req.Type {
		default:
			{
				c.WriteMessage(1, []byte("Invalid Request type"))
				fmt.Printf("Error finding request type %v", req.Type)
			}
		}
	}
}
