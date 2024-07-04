package auth

import (
	"context"
	"errors"
	"fmt"
	"todo/internal/models"
	pb "todo/pkg/proto/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AuthServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.service.Register(in.Login, in.Password)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to register user: %s, error: %s", in.Login, err.Error()))
		if errors.Is(err, models.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, models.ErrUserExists.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	s.logger.Debug(fmt.Sprintf("Sending gRPC response to register user: %s", in.Login))
	return &pb.RegisterResponse{}, nil
}

func (s *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.service.Login(in.Login, in.Password)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to login user: %s, error: %s", in.Login, err.Error()))
		if errors.Is(err, models.ErrInvalidPassword) {
			return nil, status.Error(codes.Unauthenticated, "invalid password")
		}
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	s.logger.Debug(fmt.Sprintf("Sending gRPC response with token for user: %s", in.Login))
	return &pb.LoginResponse{Token: token}, nil
}

func (s *AuthServer) UserExists(ctx context.Context, in *pb.UserExistsRequest) (*pb.UserExistsResponse, error) {
	exists, err := s.service.UserExists(in.Login)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to check if user exists: %s, error: %s", in.Login, err.Error()))
		return nil, status.Error(codes.Internal, "internal server error")
	}

	s.logger.Debug(fmt.Sprintf("Sending gRPC response to check if user exists: %s", in.Login))
	return &pb.UserExistsResponse{Exists: exists}, nil
}
