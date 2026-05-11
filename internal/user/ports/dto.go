package ports

import "time"

type CreateUserRequest struct {
	Username string `validate:"required,min=3,max=64"`
	Password string `validate:"required,min=8,max=128"`
}

type CreateUserResponse struct {
	UserID    string
	Username  string
	CreatedAt time.Time
}
