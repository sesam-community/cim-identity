package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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
	router  *httprouter.Router
	client  *http.Client
	options *serverOptions
}

// NewServer sets up and returns microservice Server
func NewServer(opt serverOptions) (*Server, error) {
	s := &Server{router: httprouter.New(), options: &opt}
	s.Routes()
	err := s.Backend()
	if err != nil {
		return nil, err
	}
	s.Logf(logLIVE, "Started RFC4122 urn:uuid-scheme UUID-v5 microservice with namespace:  %s  (\"%s\").\n", s.options.seed.String(), s.options.namespace)
	var period string
	if time.Now().Year() > 2019 {
		period = fmt.Sprintf("%d-%d", 2019, time.Now().Year())
	} else {
		period = "2019"
	}
	s.Logf(logLIVE, "Copyright Sesam.io %s. All rights reserved.\n", period)
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
