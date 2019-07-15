package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pivotal/sample-credhub-kms-plugin/plugin"
)

func main() {
	if len(os.Args) < 3  {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-unix-socket> <public-key> <private-key> \n", os.Args[0])
		os.Exit(1)
	}

	p, err := plugin.New(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	p.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	p.Stop()
}
