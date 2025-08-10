package master

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/internal/config"
	"github.com/Shariful-NomaD-Islam/ds-with-rest-grpc/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskRequest struct {
	TaskType string `json:"task_type" binding:"required"`
	Payload  string `json:"payload" binding:"required"`
}

type TaskResponse struct {
	TaskID  string `json:"task_id"`
	Success bool   `json:"success"`
	Result  string `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
}

type StatusResponse struct {
	WorkerID    string `json:"worker_id"`
	Status      string `json:"status"`
	ActiveTasks int32  `json:"active_tasks"`
}

type AllStatusResponse struct {
	Workers map[string]StatusResponse `json:"workers"`
	Total   int                       `json:"total_workers"`
}

func SetupRoutes(workerPool *WorkerPool, config *config.Config) *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now(),
			"config": gin.H{
				"server_port":   config.Server.Port,
				"workers_count": len(config.Workers),
				"grpc_timeout":  config.GRPC.Timeout,
			},
		})
	})

	// Submit task endpoint
	r.POST("/tasks", func(c *gin.Context) {
		var req TaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate task ID
		taskID := uuid.New().String()
		logger.GetLogger().Infof("Received task: %s, Type: %s, Payload: %s", taskID, req.TaskType, req.Payload)

		// Process task via worker
		resp, err := workerPool.ProcessTask(taskID, req.TaskType, req.Payload)
		if err != nil {
			logger.GetLogger().Errorf("Failed to process task %s: %v", taskID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to process task: %v", err),
			})
			return
		}

		taskResp := TaskResponse{
			TaskID:  resp.TaskId,
			Success: resp.Success,
			Result:  resp.Result,
			Error:   resp.Error,
		}

		if resp.Success {
			c.JSON(http.StatusOK, taskResp)
		} else {
			c.JSON(http.StatusBadRequest, taskResp)
		}
	})

	// Get specific worker status endpoint
	r.GET("/status/:worker_id", func(c *gin.Context) {
		workerID := c.Param("worker_id")

		resp, err := workerPool.GetWorkerStatus(workerID)
		if err != nil {
			logger.GetLogger().Errorf("Failed to get status for worker %s: %v", workerID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to get worker status: %v", err),
			})
			return
		}

		statusResp := StatusResponse{
			WorkerID:    resp.WorkerId,
			Status:      resp.Status,
			ActiveTasks: resp.ActiveTasks,
		}

		c.JSON(http.StatusOK, statusResp)
	})

	// Get all workers status endpoint
	r.GET("/status", func(c *gin.Context) {
		statuses, err := workerPool.GetAllWorkerStatuses()
		if err != nil {
			logger.GetLogger().Errorf("Failed to get all workers status: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to get workers status: %v", err),
			})
			return
		}

		workers := make(map[string]StatusResponse)
		for workerID, status := range statuses {
			workers[workerID] = StatusResponse{
				WorkerID:    status.WorkerId,
				Status:      status.Status,
				ActiveTasks: status.ActiveTasks,
			}
		}

		response := AllStatusResponse{
			Workers: workers,
			Total:   len(workers),
		}

		c.JSON(http.StatusOK, response)
	})

	return r
}
