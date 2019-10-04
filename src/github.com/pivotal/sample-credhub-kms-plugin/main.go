package main

import (
	"fmt"
	"github.com/pivotal/sample-credhub-kms-plugin-release/src/github.com/pivotal/sample-credhub-kms-plugin/plugin"
	"log"
	"os"
	"os/signal"
	"syscall"

)

func main() {
	if len(os.Args) < 2  {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-unix-socket>\n", os.Args[0])
		os.Exit(1)
	}

	p, err := plugin.New(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	p.Start()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	p.Stop()
}
