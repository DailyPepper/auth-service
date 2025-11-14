package server

import (
	"context"
	"log"

	"github.com/DailyPepper/auth-service/internal/models"
	"github.com/DailyPepper/auth-service/pkg/generated/auth"
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
