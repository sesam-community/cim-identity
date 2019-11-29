package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Options for the microservice
type Options map[string]interface{}

type serverOptions struct {
	log       io.Writer
	level     int
	seed      uuid.UUID
	namespace string
	backend   backendOptions
	defaults  defaultOptions
	options   *Options
}

type backendOptions struct {
	serviceURL   string
	transport    *http.Transport
	tokenURL     string
	clientID     string
	clientSecret string
	clientScopes []string
}

type defaultOptions struct {
	batchSize int
}

// NewOptions returns default microservice options
func NewOptions(opt *Options) serverOptions {
	var seed uuid.UUID = uuid.Nil
	var log io.Writer = os.Stdout
	levelStr := ""
	level := logERROR
	namespace := strings.Trim(os.Getenv("UUID_SEED"), " ")
	if len(namespace) == 0 {
		if opt != nil {
			// TODO: unit-tests for all of the below
			if val, exist := (*opt)["seed"]; exist && len(strings.Trim(val.(string), " ")) != 0 {
				namespace = fmt.Sprintf("%v", val)
				seed = uuid.NewSHA1(uuid.Nil, []byte(namespace))
			} else if val, exist := (*opt)["SEED"]; exist && len(strings.Trim(val.(string), " ")) != 0 {
				namespace = fmt.Sprintf("%v", val)
				seed = uuid.NewSHA1(uuid.Nil, []byte(namespace))
			} else if val, exist := (*opt)["uuid"]; exist && len(strings.Trim(val.(string), " ")) != 0 {
				seed = val.(uuid.UUID)
			} else if val, exist := (*opt)["UUID"]; exist && len(strings.Trim(val.(string), " ")) != 0 {
				seed = val.(uuid.UUID)
			}
		}
	} else {
		seed = uuid.NewSHA1(uuid.Nil, []byte(namespace))
	}
	if seed == uuid.Nil {
		fmt.Fprintf(os.Stderr, "fatal: missing environment 'UUID_SEED' or option 'seed' or 'uuid' for microservice.\n")
		time.Sleep(30 * time.Second)
		os.Exit(1)
	}

	defaults := defaultOptions{
		batchSize: 1000,
	}
	backend := backendOptions{clientScopes: []string{"data-identity"}}
	if opt != nil {
		if val, exist := (*opt)["level"]; exist {
			levelStr = fmt.Sprintf("%v", val)
		}
		if val, exist := (*opt)["log"]; exist {
			log = val.(io.Writer)
		}
		if val, exist := (*opt)["serviceURL"]; exist {
			backend.serviceURL = val.(string)
		}
		if val, exist := (*opt)["tokenURL"]; exist {
			backend.tokenURL = val.(string)
		}
		if val, exist := (*opt)["clientID"]; exist {
			backend.clientID = val.(string)
		}
		if val, exist := (*opt)["clientSecret"]; exist {
			backend.clientSecret = val.(string)
		}
		if val, exist := (*opt)["clientScopes"]; exist {
			backend.clientScopes = strings.Split(val.(string), ",")
		}
		if val, exist := (*opt)["batchSize"]; exist {
			defaults.batchSize = val.(int)
		}
	}

	if val := os.Getenv("LOG_LEVEL"); len(val) != 0 {
		levelStr = val
	}
	if len(levelStr) != 0 {
		levelStr = strings.ToUpper(levelStr)
		for k, val := range logLevelStrings {
			if val == levelStr {
				level = k
				break
			}
		}
	}

	if val := os.Getenv("SERVICE_URL"); len(val) != 0 {
		backend.serviceURL = val
	}
	if val := os.Getenv("TOKEN_URL"); len(val) != 0 {
		backend.tokenURL = val
	}
	if val := os.Getenv("CLIENT_ID"); len(val) != 0 {
		backend.clientID = val
	}
	if val := os.Getenv("CLIENT_SECRET"); len(val) != 0 {
		backend.clientSecret = val
	}
	if val := os.Getenv("CLIENT_SCOPES"); len(val) != 0 {
		backend.clientScopes = strings.Split(val, ",")
	}
	if level >= logWARN {
		if len(backend.serviceURL) == 0 {
			fmt.Fprintf(log, "missing environment 'SERVICE_URL' or option 'serviceURL' for microservice OAuth2 backend communications\n.")
		}
		if len(backend.tokenURL) == 0 {
			fmt.Fprintf(log, "missing environment 'TOKEN_URL' or option 'tokenURL' for microservice OAuth2 backend communications.\n")
		}
		if len(backend.clientID) == 0 {
			fmt.Fprintf(log, "missing environment 'CLIENT_ID' or option 'clientID' for microservice OAuth2 backend communications.\n")
		}
		if len(backend.clientSecret) == 0 {
			fmt.Fprintf(log, "missing environment 'CLIENT_SECRET' or option 'clientSecret' for microservice OAuth2 backend communications.\n")
		}
		if len(backend.clientScopes) == 0 {
			fmt.Fprintf(log, "missing environment 'CLIENT_SCOPES' or option 'clientScopes' for microservice OAuth2 backend communications.\n")
		}
	}

	if val := os.Getenv("BATCH_SIZE"); len(val) != 0 {
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			defaults.batchSize = num
		}
	}

	return serverOptions{log: log, level: level, seed: seed, namespace: namespace, backend: backend, options: opt}
}

var logLevelStrings = []string{"OFF", "CUSTOM", "QUIET", "LIVE", "FATAL", "ERROR", "WARN", "INFO", "DEBUG", "TRACE", "ALL"}

const (
	logOFF = iota
	logCUSTOM
	logQUIET
	logLIVE
	logFATAL
	logERROR
	logWARN
	logINFO
	logDEBUG
	logTRACE
	logALL
)

// Log to configured output with an INFO level
func (s *Server) Log(l string) {
	if strings.HasSuffix(l, "\n") {
		s.Logf(logINFO, l)
	} else {
		s.Logf(logINFO, "%s\n", l)
	}
}

// Logf to configured output with given level, format and parameters
func (s *Server) Logf(level int, format string, args ...interface{}) {
	if s.options.level >= level {
		fmt.Fprintf(s.options.log, format, args...)
	}
}

// Error logs to configured output with an ERROR level
func (s *Server) Error(l string) {
	if strings.HasSuffix(l, "\n") {
		s.Logf(logERROR, l)
	} else {
		s.Logf(logERROR, "%s\n", l)
	}
}

// Errorf logs to configured output with ERROR level, format and parameters
func (s *Server) Errorf(format string, args ...interface{}) {
	s.Logf(logERROR, format, args...)
}
