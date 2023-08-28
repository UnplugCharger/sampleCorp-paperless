#!/bin/bash

# Step 1: Find and kill the process running on port 8090
process_id=$(lsof -t -i:8090)
if [ -n "$process_id" ]; then
  echo "Killing existing process running on port 8090 with PID $process_id"
  kill $process_id
else
  echo "No process running on port 8090"
fi

# Give some time for the process to be killed
sleep 2

# Step 2: Rename the existing qwetu_backend_new binary if it exists
if [ -f qwetu_backend_new ]; then
    echo "Renaming existing qwetu_backend_new to qwetu_backend_old"
    mv qwetu_backend_new qwetu_backend_old
fi

# Step 3: Build the new binary
echo "Building the new binary"
go build -o qwetu_backend_new

# Step 4: Run the new binary using nohup and store logs in output.log
echo "Running the new binary with nohup"
nohup ./qwetu_backend_new > output.log 2>&1 &

echo "Service restarted successfully"
