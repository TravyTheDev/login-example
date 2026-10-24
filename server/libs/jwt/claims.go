package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewUserClaims(
	id int,
	username string,
	email string,
	role int8,
	duration time.Duration,
) (*UserClaims, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}

	return &UserClaims{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        token.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func NewImpersonatedUserClaims(
	targetID int,
	targetUsername string,
	targetEmail string,
	targetRole int8,
	impersonatedByUser ImpersonatedByUser,
	duration time.Duration,
) (*UserClaims, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}

	return &UserClaims{
		ID:             targetID,
		Username:       targetUsername,
		Email:          targetEmail,
		Role:           targetRole,
		ImpersonatedBy: &impersonatedByUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        token.String(),
			Subject:   targetEmail,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
