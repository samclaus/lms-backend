package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	w.Write([]byte("hello, world!"))

	if r.URL.Path == "/" || r.URL.Path == "/index.html" {

	} else if r.URL.Path == "/connect" {
		upgrader := websocket.Upgrader{}
		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			fmt.Println(err)
		}

		for {
			messageType, data, err := ws.ReadMessage()
			if err != nil {
				log.Printf("error reading message from websocket: %v\n", err)
				break
			}
			fmt.Println("message received: %d %s", messageType, data)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
