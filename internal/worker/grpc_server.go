package worker

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/internal/logger"
	pb "github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/pb"
)

type WorkerServer struct {
	pb.UnimplementedWorkerServiceServer
	workerID    string
	activeTasks int32
}

func NewWorkerServer(workerID string) *WorkerServer {
	return &WorkerServer{
		workerID:    workerID,
		activeTasks: 0,
	}
}

func (s *WorkerServer) ProcessTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	atomic.AddInt32(&s.activeTasks, 1)
	defer atomic.AddInt32(&s.activeTasks, -1)

	logger.GetLogger().Infof("Worker %s processing task %s of type %s", s.workerID, req.TaskId, req.TaskType)

	// Simulate task processing
	time.Sleep(2 * time.Second)

	// Simple task processing logic
	var result string
	var success bool = true
	var errorMsg string

	switch req.TaskType {
	case "compute":
		result = fmt.Sprintf("Computed result for: %s", req.Payload)
	case "process":
		result = fmt.Sprintf("Processed data: %s", req.Payload)
	default:
		success = false
		errorMsg = "Unknown task type"
		result = ""
	}

	return &pb.TaskResponse{
		TaskId:  req.TaskId,
		Success: success,
		Result:  result,
		Error:   errorMsg,
	}, nil
}

func (s *WorkerServer) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	tasks := atomic.LoadInt32(&s.activeTasks)

	return &pb.StatusResponse{
		WorkerId:    s.workerID,
		Status:      "healthy",
		ActiveTasks: tasks,
	}, nil
}
