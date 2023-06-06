package main

import (
	"context"
	"flag"
	"fmt"
	"sync"

	"client/client"
)

const (
	defaultConcurrentClientCount = 3
	defaultServerEndpointURL     = "http://localhost:8081/"
)

func parseParameters() (int, string) {
	var clientsCountFlag = flag.Int("c", defaultConcurrentClientCount, "number of concurrent client to populate.")
	var serverEndpointFlag = flag.String("u", defaultServerEndpointURL, "URL for the server endpoint.")
	flag.Parse()
	return *clientsCountFlag, *serverEndpointFlag
}

func main() {
	clientCount, serverEndpointURL := parseParameters()
	fmt.Println(clientCount)
	fmt.Println(serverEndpointURL)

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		newClient := client.NewClient(serverEndpointURL)
		go newClient.Run(ctx, wg)
	}

	fmt.Println("Press the Enter Key to stop anytime")
	fmt.Scanln()
	cancel()
	wg.Wait()
}
