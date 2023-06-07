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
	// Client.waitIntervalRandomRange is here just for exercise testing purposes in order to make it easier for the
	// tester to control, otherwise it should have been a const and not a struct field.
	waitIntervalRandomRange int
}

func (c *Client) Run(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			// For the exercise's purposes we don't care about the error received, and the response's status code is
			// printed only for the convenience of the exercise's tester, otherwise I would have ignored it as well.
			resp, _ := http.Get(c.queryURL)
			if resp != nil {
				fmt.Printf("StatusCode - %d\n", resp.StatusCode)
			}

			n := rand.Intn(c.waitIntervalRandomRange)
			time.Sleep(time.Duration(n) * time.Millisecond)
		}
	}
}

func NewClient(url string, waitIntervalRandomRange int) *Client {
	clientID := uuid.New()
	queryURL := buildQuery(url, clientID.String())
	return &Client{id: clientID,
		queryURL:                queryURL,
		waitIntervalRandomRange: waitIntervalRandomRange}
}

func buildQuery(baseURL string, clientID string) string {
	// This method is pretty specific for the given example in the exercise.
	// Decided to leave it that way to simplify things in the exercise scope, instead of making it generic and
	// supporting different types of URL formatting or receiving the full URL in NewClient.
	params := url.Values{}
	params.Add(clientIDParameterName, clientID)

	u, _ := url.ParseRequestURI(baseURL)
	u.RawQuery = params.Encode()
	return fmt.Sprintf("%v", u)
}
