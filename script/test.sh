#!/bin/bash

# Test script for the master-worker system

echo "Testing Master-Worker System"

# Start worker 1 in the background
echo "Starting worker 1..."
go run ../cmd/worker/main.go -port 50051 -id worker-1 &
WORKER1_PID=$!

# Start worker 2 in the background
echo "Starting worker 2..."
go run ../cmd/worker/main.go -port 50052 -id worker-2 &
WORKER2_PID=$!

# Wait a moment for workers to start
sleep 2

# Start master in the background
echo "Starting master..."
go run ../cmd/master/main.go -config ../config.yml &
MASTER_PID=$!

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

# Clean up - stop all processes
echo "Stopping all processes..."
kill $MASTER_PID $WORKER1_PID $WORKER2_PID

echo "Test complete!"