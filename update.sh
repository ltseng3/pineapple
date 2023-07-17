#!/bin/bash
cd ../pineapple || exit
git stash && git stash clear && git pull

go build -o program

if [ "$1" != "client" ]; then
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

  # Define the port numbers to free up
  PORT_NUMBERS=(7070 7071 7072 7078 7087)

  # Loop through the port numbers and kill the processes using each port
  for PORT_NUMBER in "${PORT_NUMBERS[@]}"; do
    PIDS=$(lsof -t -i:"$PORT_NUMBER")
    if [ -z "$PIDS" ]; then
      echo "No processes found using port $PORT_NUMBER."
    else
      echo "Killing processes using port $PORT_NUMBER..."
      kill $PIDS
    fi
  done

  if [ "$1" = "master" ]; then
    ./program -maddr "10.10.1.1" -N 3 &
    sleep 0.5
  fi

  go run server/server.go -maddr "10.10.1.1" -mport 7087 -addr "$IP" -port $PORT &
  sleep 0.5
done