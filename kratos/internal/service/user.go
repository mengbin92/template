package service

import (
	"context"
	"net/http"

	pb "explorer/api/explorer/v1"
	"explorer/internal/models/users"
	"explorer/provider/db"
	"explorer/utils"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserService struct {
	pb.UnimplementedUserServer
	UserManager *users.UserManager
	logger      *log.Helper
}

func NewUserService(userManager *users.UserManager, logger log.Logger) *UserService {
	return &UserService{
		UserManager: userManager,
		logger:      log.NewHelper(logger),
	}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	s.logger.Infof("register user: %s", req.Username)

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("register user: %s, error: %v", req.Username, err)
		return nil, err
	}

	user, err := s.UserManager.CreateUser(ctx, &users.User{
		Username: req.Username,
		Password: string(password),
		Email:    req.Email,
		Phone:    req.Phone,
	}, db.Get())
	if err != nil {
		s.logger.Errorf("register user: %s, error: %v", req.Username, err)
		return nil, err
	}

	return &pb.RegisterReply{
		UserProfile: &pb.UserProfile{
			UserId:   uint64(user.ID),
			Username: user.Username,
		},
		Token: "new user registered",
	}, nil
}
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	s.logger.Infof("login user: %s", req.Username)

	token, err := s.UserManager.Login(ctx, req.Username, req.Password, db.Get())
	if err != nil {
		s.logger.Errorf("login user: %s, error: %v", req.Username, err)
		return nil, err
	}
	return &pb.LoginReply{
		UserProfile: &pb.UserProfile{
			Username: req.Username,
		},
		Token: token,
	}, nil
}
func (s *UserService) Logout(ctx context.Context, req *emptypb.Empty) (*pb.LogoutReply, error) {
	s.logger.Info("logout user")
	ua := utils.GetUserAuthenticationFromContext(ctx)
	if ua == nil {
		s.logger.Error("logout user, user authentication not found")
		return nil, errors.New("logout user, user authentication not found")
	}
	s.UserManager.Logout(ctx, ua.UserID)
	return &pb.LogoutReply{
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "logout success",
		},
	}, nil
}
func (s *UserService) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	s.logger.Infof("update user: %s", req.Username)

	user := &users.User{
		Model: gorm.Model{
			ID: uint(req.UserId),
		},
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	err := s.UserManager.UpdateUser(ctx, user, db.Get())
	if err != nil {
		s.logger.Errorf("update user: %s, error: %v", req.Username, err)
		return nil, errors.Wrap(err, "update user failed")
	}
	return &pb.UpdateUserReply{
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "update user success",
		},
	}, nil
}
func (s *UserService) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	s.logger.Infof("delete user: %d", req.UserId)
	err := s.UserManager.DeleteUser(ctx, uint(req.UserId), db.Get())
	if err != nil {
		s.logger.Errorf("delete user: %d, error: %v", req.UserId, err)
		return nil, errors.Wrap(err, "delete user failed")
	}
	return &pb.DeleteUserReply{
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "delete user success",
		},
	}, nil
}
