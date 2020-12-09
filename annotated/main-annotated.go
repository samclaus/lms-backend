package annotated

//  import necessary packages:
//  fmt for printing to terminal
//  net/http for http server
//  gorilla from github for websockets connections
//  NOTE: must run "go get github.com/gorilla/websocket" in terminal to install gorilla
import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	//  when the program starts, listen for requests on port 8080 and handle requests with the handle() function
	http.ListenAndServe(":8080", http.HandlerFunc(handle))
}

func handle(w http.ResponseWriter, r *http.Request) { //runs every time when we receive a request
	fmt.Println(r.Method, r.URL)     //log request method and URL to terminal
	w.Write([]byte("hello, world!")) //respond to request with "hello world!" casted to a byte array

	upgrader := websocket.Upgrader{}       //create an upgrader type that will eventually upgrade the connection from http to ws
	ws, err := upgrader.Upgrade(w, r, nil) //actually attempt to upgrade the connection and create a variable to hold the error state

	//  if there is an error with upgrading the connection, print the error to the terminal
	if err != nil {
		fmt.Println(err)
	}

	ws.ReadMessage() //filler for actual server behavior which we will work on next week
}
