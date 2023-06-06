package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"server/rate_limiter"
	"time"
)

const (
	defaultServerPort           = 8081
	shutdownSeconds             = 30
	defaultLimiterSize    int64 = 5
	defaultLimiterLimit   int64 = 5
	clientIDParameterName       = "clientid"
)

var limiterSize = defaultLimiterSize
var limiterLimit = defaultLimiterLimit

type Limiter interface {
	Allow() bool
}

var clientIDToLimiter = make(map[string]Limiter)

func parseParameters() (int64, int64) {
	var limiterSizeFlag = flag.Int64("s", defaultLimiterSize, "window size for rate limit (seconds)")
	var limiterLimitFlag = flag.Int64("l", defaultLimiterLimit, "limit threshold")
	flag.Parse()
	return *limiterSizeFlag, *limiterLimitFlag
}

func getClientID(w http.ResponseWriter, r *http.Request) string {
	clientID := r.URL.Query().Get(clientIDParameterName)
	if len(clientID) == 0 {
		http.Error(w, fmt.Sprintf("'%s' parameter must be supplied.", clientIDParameterName), http.StatusUnprocessableEntity)
		return ""
	}

	_, err := uuid.Parse(clientID)
	if err != nil {
		http.Error(w, fmt.Sprintf("'%s' parameter must be valid UUID.", clientIDParameterName), http.StatusUnprocessableEntity)
		return ""
	}

	return clientID
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	clientID := getClientID(w, r)
	if clientID == "" {
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
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func main() {
	limiterSize, limiterLimit = parseParameters()

	router := http.NewServeMux()
	router.HandleFunc("/", handleRequest)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", defaultServerPort),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	fmt.Println("Press the Enter Key to stop anytime")
	fmt.Scanln()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownSeconds*time.Second)
	defer cancel()

	fmt.Printf("Shutting Down...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
