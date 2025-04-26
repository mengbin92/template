package service

import (
	"context"
	"net/http"

	pb "explorer/api/explorer/v1"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BasicService struct {
	pb.UnimplementedBasicServer
	Logger *log.Helper
}

func NewBasicService(logger log.Logger) *BasicService {
	return &BasicService{
		Logger: log.NewHelper(logger),
	}
}

func (s *BasicService) Ping(ctx context.Context, req *emptypb.Empty) (*pb.PingReply, error) {
	return &pb.PingReply{
		Status: &pb.Status{
			Code: http.StatusOK,
			Message:  "pong",
		},
	}, nil
}
