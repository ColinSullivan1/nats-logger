# A simple NATS logger
This is a simple demonstration using NATS to publish log messages in a distributed system.

## Overview
This demonstrates NATS though a simple logging utility that publishes log messages in a distributed environment.  NATS subscribers can register specific interest to filter what messages they receive.  The simple, scalable, lightweight nature of NATS makes this easy.

## Requirements
A running NATS server is required for this demonstration code.

```bash
go get github.com/nats-io/gnatsd
gnatsd
```

## Installation

```bash
go get github.com/ColinSullivan1/nats-logger
```

## API Usage
The appname is set when you create the logger, the type is set by the API called.  Let's walk through an example.

```go
import (
    "github.com/ColinSullivan1/nats-logger/natslogger"
)
```

Now create a NATS logger and log some messages.
```go
# Create a NATS logger for "myapp"
l, _ := natslogger.NewNATSLogger("myapp", "nats://localhost:4222")

l.Infof("Here is an info message.")
# This message will be published to logging.myapp.inf

l.Errorf("Here is an error message.")
# This message will be published to logging.myapp.err
```

## Subject namespace
This uses the following subject namespace.  For more on NATS subjects, see the [NATS Protocol ](http://nats.io/documentation/internals/nats-protocol/)

`logging.<appname>.<type>`
Any applications that are listening can receive messages on:

`logging.myapp.err` and `logging.myapp.inf`

...and here is the advantage of using the NATS namespace - subscribing to certain subjects using wildcards lets you filter data.

To receive all log messages from every app, subscribe to:
`logging.>`

To receive all log messages from just `myapp`, subscribe to:
`logging.myapp.*`

To receive only error messages from `myapp`, subscribe to:
`logging.myapp.err`

To receive only error messages from all applications:
`logging.*.err`

...and so forth.


## Running the demo application from the command line
```bash
Usage of demo-app:
  -app string
    	Application name (default "demoapp")
  -url string
    	URL of the nats server (default "nats://localhost:4222")
```

### Run the demo app
```bash
cd demo-app
go build
./demo-app -app myapp1 2>app1.out &
./demo-app -app myapp2 2>app2.out &
```

### Run nats-sub
Use the nats-sub example to receive and print messages being generated 
by the demo apps.

```bash
go get github.com/nats-io/go-nats
cd $GOPATH/src/github.com/nats-io/go-nats/examples
go run nats-sub.go "logging.>"
```

Example output:
```
[#1] Received on [logging.myapp2.inf]: '[myapp2] [inf] Received email avamiller@example.com from IP 126.35.98.134.'
[#2] Received on [logging.myapp1.err]: '[myapp1] [err] Error reported by benjaminanderson@example.com at IP 62.160.243.8.'
[#3] Received on [logging.myapp2.inf]: '[myapp2] [inf] Received email aubreywilliams@example.com from IP 248.51.224.20.'
[#4] Received on [logging.myapp1.err]: '[myapp1] [err] Error reported by joshuajones@test.com at IP 108.114.64.192.'
[#5] Received on [logging.myapp1.err]: '[myapp1] [err] Error reported by aidenwhite@example.com at IP 62.216.14.23.'
[#6] Received on [logging.myapp2.err]: '[myapp2] [err] Error reported by avarobinson@example.com at IP 77.160.212.140.'
[#7] Received on [logging.myapp2.inf]: '[myapp2] [inf] Received email sofiamoore@test.com from IP 172.49.237.106.'
[#8] Received on [logging.myapp1.inf]: '[myapp1] [inf] Received email avarobinson@test.com from IP 202.155.147.216.'
```

Experiment with the subscription as described above.

Enjoy!











