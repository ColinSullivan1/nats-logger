// Copyright 2017 Apcera Inc. All rights reserved.

//Package natslogger provides logging facilities to a NATS server
package natslogger

import (
	"fmt"
	"log"
	"os"

	"github.com/nats-io/go-nats"
)

const (
	// InfoLabel is the label/subject used for info log statements
	InfoLabel = "inf"

	// ErrorLabel is the label/subject used  for error log statements
	ErrorLabel = "err"

	// FatalLabel is the label/subject used  for fatal log statements
	FatalLabel = "ftl"

	// LogSubjPrefix is the subject prefix for logging.
	LogSubjPrefix = "logging"
)

// Logger is the server logger
type Logger struct {
	logger    *log.Logger // also log to the default logger
	nc        *nats.Conn  // Connection to the NATS server
	app       string      // an application name
	infoSubj  string      // subject to publish info log messages to
	errorSubj string      // subject to publish error log messages to
	fatalSubj string      // subject to publish fatal log messages to
}

// NewNATSLogger creates a logger with output directed to Stderr and as NATS messages
func NewNATSLogger(app, url string) (*Logger, error) {
	if app == "" || url == "" {
		return nil, fmt.Errorf("invalid parameter")
	}

	p := fmt.Sprintf("[%d] ", os.Getpid())
	l := &Logger{
		logger:    log.New(os.Stderr, p, log.LstdFlags),
		app:       app,
		infoSubj:  fmt.Sprintf("%s.%s.%s", LogSubjPrefix, app, InfoLabel),
		errorSubj: fmt.Sprintf("%s.%s.%s", LogSubjPrefix, app, ErrorLabel),
		fatalSubj: fmt.Sprintf("%s.%s.%s", LogSubjPrefix, app, FatalLabel),
	}

	// Create our connection to the NATS server
	var err error
	l.nc, err = nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the nats server:  %v", err)
	}
	return l, nil
}

// publish a log message to the NATS server.
func (l *Logger) publishLogMsg(subject, msg string) {
	err := l.nc.Publish(subject, []byte(msg))
	if err != nil {
		l.logger.Printf("[logging error]: couldn't publish NATS message: %v", err)
	}
}

// generate a formatted log message containing the label and application
func (l *Logger) genLogMsg(label, format string, v ...interface{}) string {
	prefix := fmt.Sprintf("[%s] [%s] ", l.app, label)
	return fmt.Sprintf(prefix+format, v...)
}

// Infof logs a notice statement
func (l *Logger) Infof(format string, v ...interface{}) {
	msg := l.genLogMsg(InfoLabel, format, v...)
	l.publishLogMsg(l.infoSubj, msg)
	l.logger.Print(msg)
}

// Errorf logs an error statement
func (l *Logger) Errorf(format string, v ...interface{}) {
	msg := l.genLogMsg(ErrorLabel, format, v...)
	l.publishLogMsg(l.errorSubj, msg)
	l.logger.Print(msg)
}

// Fatalf logs a fatal error
func (l *Logger) Fatalf(format string, v ...interface{}) {
	msg := l.genLogMsg(FatalLabel, format, v...)
	l.publishLogMsg(l.fatalSubj, msg)

	// Flush ensure the message is sent before exiting.
	l.Flush()
	l.logger.Fatalf(msg)
}

// Flush flushes the log
func (l *Logger) Flush() {
	if l.nc != nil {
		_ = l.nc.Flush() // nolint
	}
}

// Close closes the connection to NATS
func (l *Logger) Close() {
	if l.nc != nil {
		l.Flush()
		l.nc.Close()
	}
}
