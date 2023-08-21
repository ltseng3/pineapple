#!/bin/bash
cd ../pineapple || exit
git stash && git stash clear && git pull

export GOPATH=~/go/src/
export GOBIN=~/go/src/pineapple/bin
go install pineapple/src/master
go install pineapple/src/server
go install pineapple/src/client

if [ "$1" = "client" ]; then
  . test.sh
else
  HOST=$(hostname | awk -F "." '{print $1}')
  echo "Hostname: $HOST"

  if [ $HOST = "california" ]; then
      IP="10.10.1.1"
      PORT=7070
      ID=0
  elif [ $HOST = "virginia" ]; then
      IP="10.10.1.2"
      PORT=7071
      ID=1
  elif [ $HOST = "ireland" ]; then
      IP="10.10.1.3"
      PORT=7072
      ID=2
  elif [ $HOST = "oregon" ]; then
      IP="10.10.1.4"
      PORT=7073
      ID=3
  elif [ $HOST = "japan" ]; then
      IP="10.10.1.5"
      PORT=7074
      ID=4
  fi
  echo "Local IP: $IP"
  echo "Local PORT: $PORT"

  # Define the port numbers to free up
  PORT_NUMBERS=(7070 7071 7072 7073 7074 7078 7087)

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
    bin/master -maddr "10.10.1.1" -N 5 &
    sleep 0.5
  fi

  bin/server -maddr "10.10.1.1" -mport 7087 -addr "$IP" -port $PORT &
  sleep 0.5
  . test.sh $IP $PORT $ID
fi