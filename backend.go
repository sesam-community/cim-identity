package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	cred "golang.org/x/oauth2/clientcredentials"
)

/**
	* IDEAS FOR IMPROVEMENT
	*
	* The original UUID version 5 variant of this microservice use streaming
	* through io.Reader and io.Writer interfaces (http package).
	*
	* In a perfect world the backend service would support HTTP/2,
	* and streaming objects would be flowing into this microservice
	* which then would stream request objects to the backend service,
	* receive a stream of response objects from the backend resulting
	* in writing back to the Sesam origin as a stream of objects.
	* This would use as little memory as possible and be performant.
	* Instead we receive a stream in the ingress handler, bundle up for
	* a batched request to the backend service, then process the response
	* batch before writing back to the origin as part of the transform.
	* Not the ideal, but the neither the origin nor the backend seems to
	* support this streaming JSON fashion, so this batching is also saving
	* development time - for now.
**/

// Backend sets up the http.Client with authentication for backend services
func (s *Server) Backend() error {
	cfg := &s.options.backend
	if len(cfg.serviceURL) == 0 || len(cfg.clientID) == 0 || len(cfg.clientSecret) == 0 || len(cfg.tokenURL) == 0 {
		s.Logf(logWARN, "Backend services disabled due to missing configuration options. DataIdentity-API services not being used for UUIDs.\n")
		// FIXME: config-dependent UUID v5 failover or rather a logFATAL and os.Exit(1) - but what should be default?  =>  [FATAL is probably default]
		s.client = http.DefaultClient
		return nil
	}

	config := &cred.Config{
		ClientID:     cfg.clientID,
		ClientSecret: cfg.clientSecret,
		TokenURL:     cfg.tokenURL,
		Scopes:       cfg.clientScopes[:],
	}
	s.client = config.Client(context.Background())
	cfg.transport = &http.Transport{
		IdleConnTimeout: 5 * time.Minute,
	}
	return nil
}

type backendRequest struct {
	Sequence string `json:"_id"`
	Type     string `json:"type"`
	Content  string `json:"functionalId"`
}

type backendResponse struct {
	Sequence string `json:"_id"`
	Alias    string `json:"aliasId"`
	Entry    struct {
		Type    string `json:"type"`
		Content string `json:"functionalId"`
	} `json:"entity"`
	Status string `json:"status"`
}

// NewBatch creates a new Batch for transforming incoming requests into microservice response
func (s *Server) NewBatch(keyspecs []string) Batch {
	specSize := len(keyspecs)
	batchSize := s.options.defaults.batchSize
	expandedSize := batchSize * specSize
	b := Batch{
		specSize:  specSize,
		batchSize: batchSize,
		index:     -1, // is bumped very first time it will be used, so it starts at 0
		Holders:   make([]strings.Builder, expandedSize),
		Request:   make([]backendRequest, expandedSize),
		Result:    make([]map[string]interface{}, 0, batchSize),
	}
	return b
}

// SendBatch as request to the backend service, and return a response or error
func (s *Server) SendBatch(b *Batch) (resp *http.Response, err error) {
	l := len(b.Result)
	if s.options.defaults.batchSize < l {
		s.Logf(logINFO, "batch size setting increased from %d to %d  -  environment variable 'BATCH_SIZE'\n", s.options.defaults.batchSize, l)
		s.options.defaults.batchSize = l
	}
	return nil, nil
}

// WriteBatch transformed and written as response to microservice request and returns any error
func (s *Server) WriteBatch(w http.ResponseWriter, b *Batch) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return nil
}

// Batch handles batching between backend service and the resulting microservice response.
// Holders prevents the need for needing to iterate maps for allocating resolved aliases.
// Instead a simple binding is used to map multiple resolved aliases into each entity.
// This strategy approximates a flattened list of functionalIDs to aliasIDs.
type Batch struct {
	specSize  int
	batchSize int
	index     int
	Holders   []strings.Builder
	Request   []backendRequest
	Result    []map[string]interface{}
}

// increase the internal index and possibly extend internal structures if needed
func (b *Batch) increase() {
	b.index++
	expandedSize := b.batchSize * b.specSize
	if b.index >= expandedSize { // need to extend current pre-allocations
		b.Request = append(b.Request, make([]backendRequest, expandedSize)...)
		b.Holders = append(b.Holders, make([]strings.Builder, expandedSize)...)
		b.Result = append(b.Result, make([]map[string]interface{}, 0, b.batchSize)...)
		b.batchSize *= 2 // now using doubled size
	}
}

func (b *Batch) prefix(str string) *strings.Builder {
	b.increase()
	h := &b.Holders[b.index]
	l := len(str)
	if l != 0 {
		h.Grow(l + lenUUID)
		h.WriteString(str)
	} else {
		h.Grow(lenUUID)
	}
	return h
}

func (b *Batch) add(entity *map[string]interface{}) {

}
