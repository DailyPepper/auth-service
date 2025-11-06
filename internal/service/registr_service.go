package service

import (
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"context"
	"time"

	"github.com/pkg/errors"
)

type RegistrService struct {
	userRepo repository.UserRepository
}

func NewRegistrService(userRepo repository.UserRepository) *RegistrService {
	return &RegistrService{
		userRepo: userRepo,
	}
}

func (s *RegistrService) Registration(ctx context.Context, req *models.Registr) (*models.User, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check existing user")
	}
	if existingUser != nil {
		return nil, models.ErrUserAlreadyExists
	}

	// Создаем объект пользователя
	user := &models.User{
		FirstName:  req.FirstName,
		Surname:    req.Surname,
		Birthday:   req.Birthday,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   req.Password, // Пароль будет захеширован в методе
		IsActive:   true,
		IsVerified: false,
		Role:       models.RoleUser,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Хешируем пароль
	if err := user.HashPassword(); err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}

	// Сохраняем в базу
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "failed to create user in database")
	}

	// Возвращаем пользователя без пароля для безопасности
	user.Password = ""

	return user, nil
}

func (s *RegistrService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Находим пользователя по email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	if user == nil {
		return nil, models.ErrInvalidCredentials
	}

	// Проверяем активность пользователя
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// Проверяем пароль
	if !user.CheckPassword(req.Password) {
		return nil, models.ErrInvalidCredentials
	}

	// Обновляем время последнего входа
	loginTime := time.Now()
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID, loginTime); err != nil {
		return nil, errors.Wrap(err, "failed to update last login")
	}

	// Генерируем токены (заглушки - нужно реализовать JWT)
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate tokens")
	}

	// Очищаем пароль в ответе
	user.Password = ""

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *RegistrService) generateTokens(user *models.User) (string, string, error) {
	// TODO: Реализовать генерацию JWT токенов
	// Временная заглушка
	accessToken := "access_token_" + user.Email
	refreshToken := "refresh_token_" + user.Email
	return accessToken, refreshToken, nil
}

func (s *RegistrService) GetUserProfile(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Очищаем пароль
	user.Password = ""

	return user, nil
}
