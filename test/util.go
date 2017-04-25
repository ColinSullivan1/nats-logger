// Copyright 2017 Apcera Inc. All rights reserved.

package test

import (
	"fmt"
	"net"
	"time"

	"github.com/nats-io/gnatsd/logger"
	"github.com/nats-io/gnatsd/server"
)

// ClientPort is the default port for clients to connect
const ClientPort = 9876

// ClientHost is the default host for clients to connect
const ClientHost = "127.0.0.1"

// GetDefaultURL gets the default url for testing
func GetDefaultURL() string {
	return fmt.Sprintf("nats://%s:%d", ClientHost, ClientPort)
}

// RunServer runs the NATS server in a go routine
func RunServer() *server.Server {
	return RunServerWithPort(ClientPort)
}

// RunServerWithPort runs the NATS server with a monitor port in a go routine
func RunServerWithPort(cport int) *server.Server {
	var enableLogging bool

	// To enable debug/trace output in the NATS server,
	// flip the enableLogging flag.
	// enableLogging = true

	opts := &server.Options{
		Host:   ClientHost,
		Port:   cport,
		NoLog:  !enableLogging,
		NoSigs: true,
	}

	s := server.New(opts)
	if s == nil {
		panic("No NATS Server object returned.")
	}

	if enableLogging {
		l := logger.NewStdLogger(true, true, true, false, true)
		s.SetLogger(l, true, true)
	}

	// Run server in Go routine.
	go s.Start()

	end := time.Now().Add(10 * time.Second)
	for time.Now().Before(end) {
		netAddr := s.Addr()
		if netAddr == nil {
			continue
		}
		addr := s.Addr().String()
		if addr == "" {
			time.Sleep(10 * time.Millisecond)
			// Retry. We might take a little while to open a connection.
			continue
		}
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			// Retry after 50ms
			time.Sleep(50 * time.Millisecond)
			continue
		}
		_ = conn.Close() // nolint

		// Wait a bit to give a chance to the server to remove this
		// "client" from its state, which may otherwise interfere with
		// some tests.
		time.Sleep(25 * time.Millisecond)

		return s
	}
	panic("Unable to start NATS Server in Go Routine")
}
