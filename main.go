package main

import (
	"fmt"
	"net/http"
)

// This is the starting of the program
func main() {
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	w.Write([]byte("hello, world!"))
}
