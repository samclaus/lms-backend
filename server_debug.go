// +build debug

package main

import "net/http"

type debugHandler struct {
	rawHandler http.Handler
}

func (dh debugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/debug" {
		http.ServeFile(w, r, "debug/debug.html")
	} else {
		dh.rawHandler.ServeHTTP(w, r)
	}
}

func (s *Server) handler() http.Handler {
	return debugHandler{s}
}
