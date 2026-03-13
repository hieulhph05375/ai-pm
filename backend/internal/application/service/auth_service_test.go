package service

import (
	"context"
	"os"
	"project-mgmt/backend/internal/domain/entity"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

type mockUserRepo struct {
	users map[string]*entity.User
}

func (m *mockUserRepo) Create(ctx context.Context, u *entity.User) error {
	m.users[u.Email] = u
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, os.ErrNotExist
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepo) Update(ctx context.Context, u *entity.User) error {
	m.users[u.Email] = u
	return nil
}

func (m *mockUserRepo) List(ctx context.Context) ([]entity.User, error) {
	return nil, nil
}

func (m *mockUserRepo) Count(ctx context.Context, search string) (int, error) {
	return len(m.users), nil
}

func (m *mockUserRepo) ListWithPagination(ctx context.Context, search string, offset, limit int) ([]entity.User, int, error) {
	return nil, 0, nil
}

func (m *mockUserRepo) Search(ctx context.Context, query string) ([]entity.User, error) {
	return nil, nil
}

func TestAuthService_Login(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*entity.User)}
	// Note: In real test we'd use temp files for keys or a mock key provider
	// For this GSD execution, we'll assume the generated certs exist at known paths
	privPath := "../../../certs/private.pem"
	pubPath := "../../../certs/public.pem"

	s, err := NewAuthService(repo, nil, privPath, pubPath)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	user := &entity.User{
		Email: "test@example.com",
	}
	s.Register(context.Background(), user, "password")

	t.Run("Valid Credentials", func(t *testing.T) {
		accessToken, _, _, err := s.Login(context.Background(), "test@example.com", "password")
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if accessToken == "" {
			t.Error("Expected access token, got empty string")
		}

		// Verify token
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return s.(*authService).publicKey, nil
		})
		if err != nil || !token.Valid {
			t.Errorf("Token verification failed: %v", err)
		}
	})

	t.Run("Account Locking", func(t *testing.T) {
		email := "locktest@example.com"
		s.Register(context.Background(), &entity.User{Email: email}, "password")

		// Fail 5 times
		for i := 0; i < 5; i++ {
			_, _, _, err := s.Login(context.Background(), email, "wrong")
			if err == nil || err.Error() != "Email hoặc mật khẩu không đúng" {
				t.Errorf("Iteration %d: expected login failure error, got %v", i, err)
			}
		}

		// 6th attempt should be locked
		_, _, _, err := s.Login(context.Background(), email, "password")
		if err == nil {
			t.Error("Expected account to be locked, but login succeeded")
		}
		if err != nil && !contains(err.Error(), "Tài khoản bị khoá") {
			t.Errorf("Expected lock error message, got: %v", err)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))
}
