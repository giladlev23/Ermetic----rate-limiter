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
	id       uuid.UUID
}

func (c *Client) Run(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			// For the exercise's purposes we don't care about the response/error received.
			resp, err := http.Get(c.queryURL)
			if err != nil {
				fmt.Print(err)
				continue
			}
			fmt.Printf("Status Code - %d\n", resp.StatusCode)
			fmt.Printf("Status - %s\n", resp.Status)

			n := rand.Intn(5000)
			fmt.Printf("Sleeping %d milliseconds...\n\n", n)
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
	params := url.Values{}
	params.Add(clientIDParameterName, clientID)

	u, _ := url.ParseRequestURI(baseURL)
	// TODO: handle error
	u.RawQuery = params.Encode()
	return fmt.Sprintf("%v", u)
}
