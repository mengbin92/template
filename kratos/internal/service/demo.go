// Package service provides business logic services for the demo API.
package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pb "kratos-project-template/api/demo/v1"
	"kratos-project-template/internal/global"
	"kratos-project-template/provider/cache"
	"kratos-project-template/provider/db"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

// DemoService implements the demo API service.
type DemoService struct {
	pb.UnimplementedDemoServer
}

// NewDemoService creates a new instance of DemoService.
func NewDemoService() *DemoService {
	return &DemoService{}
}

// GetHello returns a hello message.
func (s *DemoService) GetHello(ctx context.Context, req *pb.GetHelloRequest) (*pb.Reply, error) {
	if req == nil {
		return s.errorReply(http.StatusBadRequest, "request cannot be nil", errors.New("nil request")), errors.New("nil request")
	}

	name := req.GetName()
	if name == "" {
		name = "World"
	}

	response := map[string]interface{}{
		"message": "Hello, " + name + "!",
		"timestamp": time.Now().Unix(),
	}

	return s.marshalAndReply(response, "marshal response error")
}

// CheckHealthy performs a health check on the service and its dependencies.
func (s *DemoService) CheckHealthy(ctx context.Context, req *pb.CheckHealthyRequest) (*pb.CheckHealthyResponse, error) {
	global.Logger.Debugf("health check requested: service=%s", req.GetService())

	serviceName := "kratos-project-template"
	if req.GetService() != "" {
		serviceName = req.GetService()
	}

	healthStatus := "healthy"
	details := make(map[string]*pb.HealthDetails)

	// Check database connection
	dbHealthy := true
	startTime := time.Now()
	if db.Get() != nil {
		sqlDB, err := db.Get().DB()
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			err := sqlDB.PingContext(pingCtx)
			latency := time.Since(startTime).Milliseconds()
			cancel()
			if err != nil {
				dbHealthy = false
				healthStatus = "degraded"
				details["database"] = &pb.HealthDetails{
					Status:    "unhealthy",
					Error:     err.Error(),
					LatencyMs: float64(latency),
				}
			} else {
				details["database"] = &pb.HealthDetails{
					Status:    "healthy",
					LatencyMs: float64(latency),
				}
			}
		} else {
			dbHealthy = false
			healthStatus = "degraded"
			details["database"] = &pb.HealthDetails{
				Status: "unhealthy",
				Error:  "failed to get database instance",
			}
		}
	} else {
		dbHealthy = false
		healthStatus = "unhealthy"
		details["database"] = &pb.HealthDetails{
			Status: "unavailable",
			Error:  "database instance is nil",
		}
	}

	// Check Redis connection (if available)
	redisStartTime := time.Now()
	redisClient := cache.GetRedisClient()
	if redisClient != nil {
		redisCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		_, err := redisClient.Ping(redisCtx).Result()
		redisLatency := time.Since(redisStartTime).Milliseconds()
		cancel()
		if err != nil {
			healthStatus = "degraded"
			details["redis"] = &pb.HealthDetails{
				Status:    "unhealthy",
				Error:     err.Error(),
				LatencyMs: float64(redisLatency),
			}
		} else {
			details["redis"] = &pb.HealthDetails{
				Status:    "healthy",
				LatencyMs: float64(redisLatency),
			}
		}
	} else {
		details["redis"] = &pb.HealthDetails{
			Status: "not_configured",
		}
	}

	// If database is critical and unhealthy, mark as unhealthy
	if !dbHealthy {
		healthStatus = "unhealthy"
	}

	response := &pb.CheckHealthyResponse{
		Status:    healthStatus,
		Service:   serviceName,
		Timestamp: time.Now().Unix(),
		Details:   details,
	}

	global.Logger.Debugf("health check completed: status=%s, service=%s", healthStatus, serviceName)
	return response, nil
}

// errorReply creates a standardized error reply.
func (s *DemoService) errorReply(code int, message string, err error) *pb.Reply {
	return &pb.Reply{
		Code:    int32(code),
		Message: fmt.Sprintf("%s: %s", message, err.Error()),
	}
}

// marshalAndReply serializes the value and returns a success reply.
func (s *DemoService) marshalAndReply(v interface{}, errorMsg string) (*pb.Reply, error) {
	jsonBytes, err := sonic.Marshal(v)
	if err != nil {
		global.Logger.Errorf("%s: %v", errorMsg, err)
		return s.errorReply(http.StatusInternalServerError, errorMsg, err), err
	}

	var structData structpb.Struct
	if err := sonic.Unmarshal(jsonBytes, &structData); err != nil {
		global.Logger.Errorf("failed to unmarshal JSON to struct: %v", err)
		return s.errorReply(http.StatusInternalServerError, "marshal error", err), err
	}

	return &pb.Reply{
		Code:    int32(http.StatusOK),
		Message: "success",
		Data:    &structData,
	}, nil
}

