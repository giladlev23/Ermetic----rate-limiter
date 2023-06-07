package main

import (
	"client/client"
	"context"
	"flag"
	"fmt"
	"sync"
)

const (
	defaultServerEndpointURL                   = "http://localhost:8081/"
	defaultConcurrentClientCount               = 1
	defaultWaitIntervalRandomRangeMilliseconds = 1000
)

type Params struct {
	clientsCount                        int
	serverEndpointURL                   string
	waitIntervalRandomRangeMilliseconds int
}

func parseParameters() *Params {
	var clientsCountFlag = flag.Int("c", defaultConcurrentClientCount, "number of concurrent client to populate.")
	var serverEndpointFlag = flag.String("u", defaultServerEndpointURL, "URL for the server endpoint.")
	var waitIntervalRandomRangeMillisecondsFlag = flag.Int("w", defaultWaitIntervalRandomRangeMilliseconds, "Wait interval between each client requests as random range (milliseconds).")
	flag.Parse()

	return &Params{*clientsCountFlag, *serverEndpointFlag, *waitIntervalRandomRangeMillisecondsFlag}
}

func main() {
	params := parseParameters()

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < params.clientsCount; i++ {
		wg.Add(1)
		newClient := client.NewClient(params.serverEndpointURL, params.waitIntervalRandomRangeMilliseconds)
		go newClient.Run(ctx, wg)
	}

	fmt.Println("Press the Enter Key to stop anytime")
	fmt.Scanln()
	cancel()
	wg.Wait()
}
