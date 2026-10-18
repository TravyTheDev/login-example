package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     int8   `json:"role"`
	jwt.RegisteredClaims
}

type PasswordResetClaims struct {
	Email   string `json:"email"`
	Numbers int    `json:"numbers"`
	jwt.RegisteredClaims
}

type Session struct {
	ID           string    `json:"id"`
	UserEmail    string    `json:"userEmail"`
	RefreshToken string    `json:"refreshToken"`
	IsRevoked    bool      `json:"isRevoked"`
	CreatedAt    time.Time `json:"createdAt"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
