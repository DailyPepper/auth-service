package models

import "time"

type Registr struct {
	FirstName string    `json:"first_name" validate:"required,min=2,max=50"`
	Surname   string    `json:"surname" validate:"required,min=2,max=50"`
	Birthday  time.Time `json:"birthday" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8"`
	Phone     *string   `json:"phone,omitempty" validate:"omitempty,e164"`
}
