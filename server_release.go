// +build !debug

package main

import "net/http"

func (s *Server) handler() http.Handler {
	return s
}
