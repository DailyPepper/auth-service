package server

import (
	"context"
	"log"
	"time"

	"github.com/DailyPepper/auth-service/internal/models"
	"github.com/DailyPepper/auth-service/pkg/generated/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	log.Printf("gRPC Register called for email: %s", req.Email)

	registrModel := &models.Registr{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		Surname:   req.Surname,
	}

	user, err := s.registrService.Registration(ctx, registrModel)
	if err != nil {
		return nil, s.mapErrorToStatus(err)
	}

	return &auth.RegisterResponse{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		Surname:   user.Surname,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	log.Printf("gRPC Login called for email: %s", req.Email)

	loginModel := &models.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	loginResponse, err := s.registrService.Login(ctx, loginModel)
	if err != nil {
		return nil, s.mapErrorToStatus(err)
	}

	return &auth.LoginResponse{
		AccessToken:  loginResponse.AccessToken,
		RefreshToken: loginResponse.RefreshToken,
		ExpiresAt:    timestamppb.New(time.Now().Add(24 * time.Hour)), // TODO: использовать реальное время истечения
	}, nil
}

func (s *GRPCServer) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	log.Printf("gRPC ValidateToken called")

	user, err := s.registrService.ValidateToken(ctx, req.Token)
	if err != nil {
		return &auth.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &auth.ValidateTokenResponse{
		Valid:  true,
		UserId: string(rune(user.ID)),
		Email:  user.Email,
	}, nil
}

func (s *GRPCServer) mapErrorToStatus(err error) error {
	switch err {
	case models.ErrUserAlreadyExists:
		return status.Error(codes.AlreadyExists, "user with this email already exists")
	case models.ErrInvalidCredentials:
		return status.Error(codes.Unauthenticated, "invalid email or password")
	case models.ErrUserNotFound:
		return status.Error(codes.NotFound, "user not found")
	default:
		log.Printf("Internal error: %v", err)
		return status.Error(codes.Internal, "internal server error")
	}
}
