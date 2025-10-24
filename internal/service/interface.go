package service

import (
	"auth-service/internal/models"
	"context"
)

type Registr interface {
	Registration(ctx context.Context, req *models.Registr) (*models.User, error)
}
