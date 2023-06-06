package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"server/rate_limiter"
	"time"
)

const (
	defaultLimiterSize    int64 = 5
	defaultLimiterLimit   int64 = 2
	clientIDParameterName       = "clientid"
)

var limiterSize = defaultLimiterSize
var limiterLimit = defaultLimiterLimit

type Limiter interface {
	Allow() bool
}

var clientIDToLimiter = make(map[string]Limiter)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	clientID := r.URL.Query().Get(clientIDParameterName)
	// TODO: Add input validation
	if len(clientID) == 0 {
		http.Error(w, fmt.Sprintf("'%s' parameter must be supplied.", clientIDParameterName), http.StatusUnprocessableEntity)
		return
	}

	allowed := false

	if limiter, ok := clientIDToLimiter[clientID]; ok {
		allowed = limiter.Allow()
	} else {
		newLimiter := rate_limiter.NewLimiter(time.Duration(limiterSize)*time.Second, limiterLimit)
		clientIDToLimiter[clientID] = newLimiter
		allowed = newLimiter.Allow()
	}

	if allowed {
		//w.WriteHeader(http.StatusOK)
		http.Error(w, "200!", http.StatusOK)
	} else {
		//w.WriteHeader(http.StatusServiceUnavailable)
		http.Error(w, "503!", http.StatusServiceUnavailable)
	}
}

func parseParameters() (int64, int64) {
	var limiterSizeFlag = flag.Int64("s", defaultLimiterSize, "window size for rate limit (seconds)")
	var limiterLimitFlag = flag.Int64("l", defaultLimiterLimit, "limit threshold")
	flag.Parse()
	return *limiterSizeFlag, *limiterLimitFlag
}

func main() {
	limiterSize, limiterLimit = parseParameters()

	// TODO: Add concurrency
	http.HandleFunc("/", handleRequest)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
