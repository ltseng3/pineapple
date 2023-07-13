#!/bin/bash
cd ../pineapple || exit
git stash && git stash clear && git pull

HOST=$(hostname | awk -F "." '{print $1}')
echo "Hostname: $HOST"

if [ $HOST = "node-0" ]; then
    IP="10.10.1.1"
    PORT=7070
elif [ $HOST = "node-1" ]; then
    IP="10.10.1.2"
    PORT=7071
elif [ $HOST = "node-2" ]; then
    IP="10.10.1.3"
    PORT=7072
fi
echo "Local IP: $IP"
echo "Local PORT: $PORT"

go build

if [ "$1" = "master" ]; then
  ./pineapple -maddr "10.10.1.1" -N 3 &
  sleep 0.5
fi

go run server/server.go -maddr "10.10.1.1" -mport 7087 -addr "$IP" -port $PORT &
sleep 0.5