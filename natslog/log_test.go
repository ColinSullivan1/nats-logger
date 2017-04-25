// Copyright 2017 Apcera Inc. All rights reserved.

package natslog

import (
	"strings"
	"testing"
	"time"

	lt "github.com/ColinSullivan1/nats-logger/test"
	nats "github.com/nats-io/go-nats"
)

// Test the logger when the server is not running
func TestLoggerNoServer(t *testing.T) {
	_, err := NewNATSLogger("myapp", lt.GetDefaultURL())
	if err == nil {
		t.Fatal("did not receive expected error")
	}
}

// Convenience function to check the content of the next message
func checkNextMsg(t *testing.T, sub *nats.Subscription, subject string, expected ...string) {
	var err error

	// Get the next message from the NATS server
	m, err := sub.NextMsg(time.Second * 5)
	if err != nil {
		t.Fatalf("error reading message: %v", err)
	}

	// Check that the subject is what we expect
	if m.Subject != subject {
		t.Fatalf("NATS log message subject %s does not equal expected subject %q", m.Subject, subject)
	}

	// Check the data contents for expected values
	val := string(m.Data)
	for _, r := range expected {
		if !strings.Contains(val, r) {
			t.Fatalf("NATS log message %q does not contain %q", val, r)
		}
	}
}

// creates a NATS connection and subscription
func createNATSConnAndSub(t *testing.T) (*nats.Conn, *nats.Subscription) {
	nc, err := nats.Connect(lt.GetDefaultURL())
	if err != nil {
		t.Fatalf("couldn't connect to the NATS server: %v", err)
	}

	sub, err := nc.SubscribeSync(">")
	if err != nil {
		t.Fatalf("couldn't subscribe: %v", err)
	}

	_ = nc.Flush()

	return nc, sub
}

// TestLoggerBasic some basic testing.
func TestLogger(t *testing.T) {
	s := lt.RunServer()
	defer s.Shutdown()

	l, err := NewNATSLogger("myapp", lt.GetDefaultURL())
	if err != nil {
		t.Fatal("did not receive expected error")
	}

	nc, sub := createNATSConnAndSub(t)
	defer nc.Close()

	l.Infof("Notice Test")
	checkNextMsg(t, sub, "logging.myapp.inf", "Notice Test", InfoLabel, "myapp")

	l.Errorf("Error Test")
	checkNextMsg(t, sub, "logging.myapp.err", "Error Test", ErrorLabel, "myapp")

	// check that flush doesn't panic
	l.Flush()

	// close the logger
	l.Close()
}

// TestLoggerBasic some basic testing.
func TestLoggerLostConnection(t *testing.T) {
	s := lt.RunServer()
	defer s.Shutdown()

	var err error
	l, err := NewNATSLogger("myapp", lt.GetDefaultURL())
	if err != nil {
		t.Fatal("did not receive expected error")
	}

	s.Shutdown()

	// make sure there's no panic, etc.
	l.Infof("Notice Test")
}

// TestLoggerInvalidSetup test bad params, no server, etc.
func TestLoggerBadParameters(t *testing.T) {

	// no server
	if _, err := NewNATSLogger("myapp", "nats://127.0.0.1:10456"); err == nil {
		t.Fatal("did not receive expected error")
	}

	// invalid app
	if _, err := NewNATSLogger("", lt.GetDefaultURL()); err == nil {
		t.Fatal("did not receive expected error")
	}

	// invalid url
	if _, err := NewNATSLogger("myapp", ""); err == nil {
		t.Fatal("did not receive expected error")
	}
}

// TestLoggerFatalf checks for a panic
func TestLoggerFatalf(t *testing.T) {

	//TODO - FIXME.
	s := lt.RunServer()
	defer s.Shutdown()

	l, err := NewNATSLogger("myapp", lt.GetDefaultURL())
	if err != nil {
		t.Fatal("did not receive expected error")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("log fatalf should panic")
		}
	}()

	l.Fatalf("Fatal Test")
}
