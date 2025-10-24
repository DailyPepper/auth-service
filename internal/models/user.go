package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name" validate:"required,min=2,max=50"`
	Surname   string    `json:"surname" db:"surname" validate:"required,min=2,max=50"`
	Birthday  time.Time `json:"birthday" db:"birthday" validate:"required"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Password  string    `json:"-" db:"password_hash"`
	Phone     *string   `json:"phone,omitempty" db:"phone" validate:"omitempty,e164"`

	// Статус и временные метки
	IsActive   bool       `json:"is_active" db:"is_active"`
	IsVerified bool       `json:"is_verified" db:"is_verified"`
	LastLogin  *time.Time `json:"last_login,omitempty" db:"last_login"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`

	Role UserRole `json:"role" db:"role"`
}

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Метод для возврата профиля (без чувствительных данных)
func (u *User) ToProfile() *User {
	return u
}
