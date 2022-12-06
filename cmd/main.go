package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/deanfoley/netspeak-go-test/internal/api"
)

func main() {
	readyCheck := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	apiDone := make(chan struct{}, 1)
	closedServices := make(chan struct{})

	go api.Launch(readyCheck, apiDone, closedServices)
	<-readyCheck

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		close(apiDone)
	}()

	<-closedServices

	fmt.Println("Server shut down. Goodbye!")
}
