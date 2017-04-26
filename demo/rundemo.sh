#!/bin/sh

myip=`ifconfig | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p' | grep 10.0`

docker build -t log-demo-app .

echo "Launching NATS Server in a container."
docker run -d --rm -p4222:4222 --name nats-server nats 

echo "Launching demo app myapp1 in a container."
docker run -d --rm --name myapp1 log-demo-app -url nats://$myip:4222 -app myapp1

echo "Launching demo app myapp1 in a container."
docker run -d --rm --name myapp2 log-demo-app -url nats://$myip:4222 -app myapp2

echo "To listen to published log messages, run:"
echo "nats-sub -s nats://$myip:4222 \">\""

