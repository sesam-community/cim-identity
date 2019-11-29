package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s, err := NewServer(NewOptions(nil))
	if err != nil {
		return err
	}
	return http.ListenAndServe(":5000", s)
}

// Server is a simple microservice
type Server struct {
	router    *httprouter.Router
	client    *http.Client
	transport *http.Transport
	options   *serverOptions
}

// NewServer sets up and returns microservice Server
func NewServer(opt serverOptions) (*Server, error) {
	s := &Server{router: httprouter.New(), options: &opt}
	s.Routes()
	err := s.Backend()
	if err != nil {
		return nil, err
	}
	s.Logf(logLIVE, "Started Hafslund Nett RFC4122 urn:uuid-scheme UUID-v4 microservice.\n")
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
