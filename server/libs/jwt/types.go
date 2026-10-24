package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type ImpersonatedByUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserClaims struct {
	ID             int                 `json:"id"`
	Username       string              `json:"username"`
	Email          string              `json:"email"`
	Role           int8                `json:"role"`
	ImpersonatedBy *ImpersonatedByUser `json:"impersonatedBy,omitempty"`
	jwt.RegisteredClaims
}

type PasswordResetClaims struct {
	Email   string `json:"email"`
	Numbers int    `json:"numbers"`
	jwt.RegisteredClaims
}
