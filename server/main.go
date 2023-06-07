package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"server/rate_limiter"
	"server/utils"
	"sync"
	"time"
)

const (
	shutdownSeconds   = 30
	defaultServerPort = 8081

	defaultRateLimit         int64 = 5
	defaultWindowSecondsSize int64 = 5

	clientIDParameterName = "clientid"
)

type Limiter interface {
	Allow() bool
}

func newRateLimitRequestHandler(p *params) *rateLimitRequestHandler {
	return &rateLimitRequestHandler{rateLimit: p.rateLimit,
		windowSizeSeconds: time.Duration(p.windowSizeSeconds) * time.Second,
		lock:              sync.RWMutex{},
		clientIDToLimiter: make(map[string]Limiter),
	}
}

type rateLimitRequestHandler struct {
	rateLimit         int64
	windowSizeSeconds time.Duration

	lock              sync.RWMutex
	clientIDToLimiter map[string]Limiter
}

func (h *rateLimitRequestHandler) isAllowed(clientID string) bool {
	h.lock.RLock()
	limiter, ok := h.clientIDToLimiter[clientID]
	h.lock.RUnlock()

	if ok {
		return limiter.Allow()
	} else {
		newLimiter := rate_limiter.NewLimiter(h.windowSizeSeconds, h.rateLimit)
		h.lock.Lock()
		h.clientIDToLimiter[clientID] = newLimiter
		h.lock.Unlock()

		return newLimiter.Allow()
	}
}

func (h *rateLimitRequestHandler) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// Decided to not distinguish between HTTP methods in rate limiting.
		// Note - in the current solution a single clientID can send over-the-limit requests to a different URL path
		// and get 404 instead of 503 - this is by choice for exercise simplicity purposes.
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	clientID, err := utils.ParseIDParameter(r, clientIDParameterName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if h.isAllowed(clientID) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

type params struct {
	rateLimit         int64
	windowSizeSeconds int64
}

func parseParameters() *params {
	var rateLimitFlag = flag.Int64("r", defaultRateLimit, "rate limit")
	var windowSizeSecondsFlag = flag.Int64("s", defaultWindowSecondsSize, "window size for rate limit (seconds)")
	flag.Parse()

	return &params{*rateLimitFlag, *windowSizeSecondsFlag}
}

func main() {
	parameters := parseParameters()
	handler := newRateLimitRequestHandler(parameters)

	router := http.NewServeMux()
	router.HandleFunc("/", handler.handleRequest)

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
