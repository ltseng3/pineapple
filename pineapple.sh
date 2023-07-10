#Run from root directory

#!/bin/bash

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

# After cleaning, run pineapple.
go build
./pineapple -N 3 &
sleep 0.5

go run server/server.go -port 7070 &
sleep 0.5

go run server/server.go -port 7071 &
sleep 0.5

go run server/server.go -port 7072 &
sleep 0.5