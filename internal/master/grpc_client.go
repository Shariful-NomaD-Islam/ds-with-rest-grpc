package master

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"master-worker-system/internal/config"
	pb "master-worker-system/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerClient struct {
	conn    *grpc.ClientConn
	client  pb.WorkerServiceClient
	addr    string
	id      string
	timeout time.Duration
}

type WorkerPool struct {
	workers []*WorkerClient
	config  *config.Config
	counter atomic.Int64 
}

func NewWorkerPool(config *config.Config) (*WorkerPool, error) {
	pool := &WorkerPool{
		workers: make([]*WorkerClient, 0, len(config.Workers)),
		config:  config,
	}

	timeout := config.GetGRPCTimeout()

	for _, workerConfig := range config.Workers {
		log.Printf("Connecting to worker %s at %s", workerConfig.ID, workerConfig.URL)

		conn, err := grpc.Dial(
			workerConfig.URL,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithTimeout(timeout),
		)
		if err != nil {
			log.Printf("Failed to connect to worker %s at %s: %v", workerConfig.ID, workerConfig.URL, err)
			continue
		}

		client := pb.NewWorkerServiceClient(conn)
		pool.workers = append(pool.workers, &WorkerClient{
			conn:    conn,
			client:  client,
			addr:    workerConfig.URL,
			id:      workerConfig.ID,
			timeout: timeout,
		})
	}

	if len(pool.workers) == 0 {
		return nil, fmt.Errorf("failed to connect to any workers")
	}

	return pool, nil
}

func (p *WorkerPool) ProcessTask(taskID, taskType, payload string) (*pb.TaskResponse, error) {
	if len(p.workers) == 0 {
		return nil, fmt.Errorf("no workers available")
	}

	// Simple round-robin selection
    workerIdx := p.counter.Add(1) % int64(len(p.workers))

    // Optional: reset the counter when it wraps around
    if workerIdx == 0 {
        p.counter.Store(0)
    }

    worker := p.workers[workerIdx]

	ctx, cancel := context.WithTimeout(context.Background(), worker.timeout)
	defer cancel()

	req := &pb.TaskRequest{
		TaskId:   taskID,
		TaskType: taskType,
		Payload:  payload,
	}

	return worker.client.ProcessTask(ctx, req)
}

func (p *WorkerPool) GetWorkerStatus(workerID string) (*pb.StatusResponse, error) {
	for _, worker := range p.workers {
		if worker.id == workerID {
			ctx, cancel := context.WithTimeout(context.Background(), worker.timeout)
			defer cancel()

			req := &pb.StatusRequest{
				WorkerId: workerID,
			}

			return worker.client.GetStatus(ctx, req)
		}
	}

	return nil, fmt.Errorf("worker %s not found", workerID)
}

func (p *WorkerPool) GetAllWorkerStatuses() (map[string]*pb.StatusResponse, error) {
	statuses := make(map[string]*pb.StatusResponse)

	for _, worker := range p.workers {
		ctx, cancel := context.WithTimeout(context.Background(), worker.timeout)
		defer cancel()

		req := &pb.StatusRequest{
			WorkerId: worker.id,
		}

		status, err := worker.client.GetStatus(ctx, req)
		if err != nil {
			log.Printf("Failed to get status for worker %s: %v", worker.id, err)
			continue
		}

		statuses[worker.id] = status
	}

	return statuses, nil
}

func (p *WorkerPool) Close() {
	for _, worker := range p.workers {
		if worker.conn != nil {
			worker.conn.Close()
		}
	}
}
