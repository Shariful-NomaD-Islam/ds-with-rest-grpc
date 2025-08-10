#!/bin/bash

# Test script for the master-worker system

echo "Testing Master-Worker System"

# Function to clean up processes
cleanup() {
    echo "Stopping all processes..."
    pkill -f "master-worker-system/cmd/master"
    pkill -f "master-worker-system/cmd/worker"
    echo "Cleanup complete."
}

# Trap signals to ensure cleanup on exit
trap cleanup EXIT

# Ensure no old processes are running
cleanup

# Start worker 1 in the background
echo "Starting worker 1..."
go run ../cmd/worker/main.go -port 50051 -id worker-1 &

# Start worker 2 in the background
echo "Starting worker 2..."
go run ../cmd/worker/main.go -port 50052 -id worker-2 &

# Wait a moment for workers to start
sleep 2

# Start master in the background
echo "Starting master..."
go run ../cmd/master/main.go -config ../config.yml &

# Wait for master to start
sleep 3

# Test health endpoint
echo "Testing health endpoint..."
curl -s http://localhost:8080/health | jq .

# Test submitting a task
echo "Submitting a task..."
curl -s -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"task_type": "compute", "payload": "2+2"}' | jq .

# Test getting all worker statuses
echo "Getting all worker statuses..."
curl -s http://localhost:8080/status | jq .

# Test getting specific worker status
echo "Getting specific worker status..."
curl -s http://localhost:8080/status/worker-1 | jq .

echo "Test complete!"
