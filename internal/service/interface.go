package service

import (
	"context"

	"github.com/DailyPepper/auth-service/internal/models"
)

type Registr interface {
	Registration(ctx context.Context, req *models.Registr) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	GetUserProfile(ctx context.Context, userID int64) (*models.User, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}
