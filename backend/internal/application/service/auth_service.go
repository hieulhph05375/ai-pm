package service

import (
	"context"
	"crypto/rsa"
	"project-mgmt/backend/internal/domain/entity"
)

type AuthService interface {
	Register(ctx context.Context, user *entity.User, password string) error
	Login(ctx context.Context, email, password string) (accessToken string, refreshToken string, user *entity.User, err error)
	Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error)
	GetPublicKey() *rsa.PublicKey
}
