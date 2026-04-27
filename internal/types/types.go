package types

import (
	"college-diary/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	Role models.Role `json:"role"`
	jwt.RegisteredClaims
}