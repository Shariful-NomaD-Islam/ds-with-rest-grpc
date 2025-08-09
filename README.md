# Master-Worker System

A distributed task processing system with a master node that distributes tasks to worker nodes.

## Project Structure

This project follows Go idiomatic conventions:

```
.
├── build
│   ├── master           # Binary for running master node
│   └── worker           # Binary for running worker node(s)
├── cmd/                 # Main applications
│   ├── master/          # Master node main package
│   │   └── main.go
│   └── worker/          # Worker node main package
│       └── main.go
├── internal/            # Private application code
│   ├── config/          # Configuration handling
│   │   └── config.go
│   ├── master/          # Master business logic
│   │   ├── grpc_client.go
│   │   └── handlers.go
│   └── worker/          # Worker business logic
│       └── grpc_server.go
├── pb/                  # Generated protobuf code
├── proto/               # Protocol buffer definitions
├── script
│   └── test.sh          # Shell script for testing the system
├── config.yml           # Configuration file
├── go.mod               # Go module definition
└── go.sum               # Go module checksums
```

## Key Features

- Master node distributes tasks to worker nodes via gRPC
- RESTful HTTP API for submitting tasks and checking status
- YAML-based configuration
- Graceful shutdown handling
- Health monitoring of workers

## Getting Started

### Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)

### Building

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Generate protobuf code:
   ```bash
   protoc --go_out=. --go-grpc_out=. proto/worker.proto
   ```

3. Build the master:
   ```bash
   go build -o build/master cmd/master/main.go
   ```

4. Build the worker:
   ```bash
   go build -o build/worker cmd/worker/main.go
   ```

### Running

1. Start one or more workers:
   ```bash
   cd build
   ./worker -port 50051 -id worker-1
   ./worker -port 50052 -id worker-2
   ```

2. Start the master:
   ```bash
   cd build
   ./master -config ../config.yml
   ```

### API Endpoints

- `GET /health` - Health check
- `POST /tasks` - Submit a task
- `GET /status` - Get status of all workers
- `GET /status/:worker_id` - Get status of a specific worker

### Configuration

The master node is configured via `config.yml`:

```yaml
server:
  port: "8080"
  host: "localhost"

workers:
  - url: "localhost:50051"
    id: "worker-1"
  - url: "localhost:50052"
    id: "worker-2"

grpc:
  timeout: "10s"
  max_retries: 3

logging:
  level: "warn"
```

### Test
```bash
   cd script 
   ./test.sh || echo "Test script failed"
```