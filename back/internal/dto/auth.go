package dto

import "pouic/internal/model"

type RegisterDto struct {
	Name string     `json:"name" validate:"required, min=3, max=20"`
	Password string `json:"password" validate:"required, min=2, max=72"`
}

type LoginDto struct {
	RegisterDto
}

type AuthResponse struct {
	Token string        `json:"token"`
	Player model.Player `json:"player"`
}