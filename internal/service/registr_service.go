package service

import (
	"auth-service/internal/models"
	"context"
	"time"
)

type RegistrService struct {
	// зависимости
}

func NewRegistrService() *RegistrService {
	return &RegistrService{}
}

func (s *RegistrService) Registration(ctx context.Context, req *models.Registr) (*models.User, error) {
	// Заглушка для тестирования
	return &models.User{
		ID:        1,
		FirstName: req.FirstName,
		Surname:   req.Surname,
		Birthday:  req.Birthday,
		Email:     req.Email,
		Phone:     req.Phone,
		CreatedAt: time.Now(),
	}, nil
}
