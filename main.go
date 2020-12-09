package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handle))
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	w.Write([]byte("hello, world!"))
}
