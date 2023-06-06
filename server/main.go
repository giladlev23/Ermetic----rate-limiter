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
	"sync"
	"time"
)

const (
	shutdownSeconds   = 30
	defaultServerPort = 8081

	defaultRateLimit  int64 = 5
	defaultWindowSize int64 = 5

	clientIDParameterName = "clientid"
)

var params Params
var clientIDToLimiter = make(map[string]Limiter)
var lock = sync.RWMutex{}

type Limiter interface {
	Allow() bool
}

type Params struct {
	rateLimit  int64
	windowSize int64
}

func parseParameters() {
	var rateLimitFlag = flag.Int64("r", defaultRateLimit, "rate limit")
	var windowSizeFlag = flag.Int64("s", defaultWindowSize, "window size for rate limit (seconds)")
	flag.Parse()

	params = Params{*rateLimitFlag, *windowSizeFlag}
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

	lock.RLock()
	limiter, ok := clientIDToLimiter[clientID]
	lock.RUnlock()

	if ok {
		allowed = limiter.Allow()
	} else {
		newLimiter := rate_limiter.NewLimiter(time.Duration(params.windowSize)*time.Second, params.rateLimit)
		lock.Lock()
		clientIDToLimiter[clientID] = newLimiter
		lock.Unlock()
		allowed = newLimiter.Allow()
	}

	if allowed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func main() {
	parseParameters()

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
