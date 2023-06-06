package client

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/url"
	"sync"

	"net/http"
	"time"
)

const clientIDParameterName = "clientid"

type Client struct {
	queryURL string
	// Client.id has no purpose, but it seems right for me to have the ID available for future usage.
	id uuid.UUID
}

func (c *Client) Run(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			// For the exercise's purposes we don't care about the response/error received.
			_, _ = http.Get(c.queryURL)

			n := rand.Intn(1000)
			time.Sleep(time.Duration(n) * time.Millisecond)
		}
	}
}

func NewClient(url string) *Client {
	clientID := uuid.New()
	queryURL := buildQuery(url, clientID.String())
	return &Client{id: clientID,
		queryURL: queryURL}
}

func buildQuery(baseURL string, clientID string) string {
	// This method is pretty specific for the given example in the exercise.
	// Decided to leave it that way to simplify things in the exercise scope, instead of making it generic and
	// receiving the full URL in NewClient.
	params := url.Values{}
	params.Add(clientIDParameterName, clientID)

	u, _ := url.ParseRequestURI(baseURL)
	u.RawQuery = params.Encode()
	return fmt.Sprintf("%v", u)
}
