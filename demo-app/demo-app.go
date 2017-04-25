package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/ColinSullivan1/nats-logger/natslogger"
	"github.com/Pallinder/go-randomdata"
	"github.com/nats-io/go-nats"
)

const (
	defaultAppName = "demoapp"
)

// Randomly generates log statements, about 1/3 are error statements
func generateLogStatements(l *natslogger.Logger) {
	for {
		if rand.Intn(2) > 0 {
			l.Infof(fmt.Sprintf("Received email %s from IP %s.",
				randomdata.Email(), randomdata.IpV4Address()))
		} else {
			l.Errorf(fmt.Sprintf("Error reported by %s at IP %s.",
				randomdata.Email(), randomdata.IpV4Address()))
		}

		time.Sleep(time.Duration(rand.Intn(450)+50) * time.Millisecond)
	}
}

func main() {
	var appName string
	var nURL string

	// Parse flags
	flag.StringVar(&appName, "app", defaultAppName, "Application name")
	flag.StringVar(&nURL, "url", nats.DefaultURL, "URL of the NATS server")

	// get the app name and the NATS url
	flag.Parse()

	// create a new NATS logger
	l, err := natslogger.NewNATSLogger(appName, nURL)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(-1)
	}

	// shut down on interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		// cleanly close the logger
		l.Close()
		os.Exit(0)
	}()

	// Now go generate some log statements.
	go generateLogStatements(l)

	runtime.Goexit()
}
