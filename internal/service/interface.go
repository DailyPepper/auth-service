package service

import (
	"context"

	"github.com/DailyPepper/auth-service/internal/models"
)

type Registr interface {
	Registration(ctx context.Context, req *models.Registr) (*models.User, error)
}
