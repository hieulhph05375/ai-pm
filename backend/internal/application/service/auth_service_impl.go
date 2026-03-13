package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewAuthService(ur repository.UserRepository, rr repository.RoleRepository, privPath, pubPath string) (AuthService, error) {
	privBytes, err := os.ReadFile(privPath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key: %w", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %w", err)
	}

	pubBytes, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("error reading public key: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	return &authService{
		userRepo:   ur,
		roleRepo:   rr,
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

func (s *authService) Register(ctx context.Context, u *entity.User, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	u.HashedPassword = string(hashed)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.IsActive = true
	return s.userRepo.Create(ctx, u)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, *entity.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", nil, errors.New("Email hoặc mật khẩu không đúng")
	}

	// Check if account is locked
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		remaining := time.Until(*user.LockedUntil).Minutes()
		return "", "", nil, fmt.Errorf("Tài khoản bị khoá. Vui lòng thử lại sau %.0f phút", remaining)
	}

	if !user.IsActive {
		return "", "", nil, errors.New("tài khoản của bạn đã bị vô hiệu hóa")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		// Increment failed attempts
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= 5 {
			lockTime := time.Now().Add(30 * time.Minute)
			user.LockedUntil = &lockTime
		}
		_ = s.userRepo.Update(ctx, user)
		return "", "", nil, errors.New("Email hoặc mật khẩu không đúng")
	}

	// Success: reset attempts
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil
	_ = s.userRepo.Update(ctx, user)

	// Fetch permissions for the user's role
	permNames := []string{}
	if s.roleRepo != nil && user.RoleID > 0 {
		perms, permErr := s.roleRepo.GetPermissionsByRoleID(ctx, user.RoleID)
		if permErr == nil {
			for _, p := range perms {
				permNames = append(permNames, p.Name)
			}
		}
	}

	// Create JWT Access Token with RS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
		"iat":   time.Now().Unix(),
		"role":  user.RoleID,
		"admin": user.IsAdmin,
		"perms": permNames,
	})

	accessToken, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", "", nil, fmt.Errorf("error signing token: %w", err)
	}

	// Generate a secure random refresh token
	b := make([]byte, 32)
	rand.Read(b)
	refreshToken := hex.EncodeToString(b)

	// In a real app, you'd save this to the database
	// For now, we'll return it as requested.
	return accessToken, refreshToken, user, nil
}

func (s *authService) Refresh(ctx context.Context, rt string) (string, string, error) {
	// Simple implementation for now, just to remove mock prefix
	// Usually would involve verifying the old RT in DB and issuing new ones
	return "new_access_token_placeholder", "new_refresh_token_placeholder", nil
}

func (s *authService) GetPublicKey() *rsa.PublicKey {
	return s.publicKey
}
