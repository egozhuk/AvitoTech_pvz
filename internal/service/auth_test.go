package service_test

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"pvs/internal/domain"
	"pvs/internal/service"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if u := args.Get(0); u != nil {
		return u.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

var jwtSecret = []byte("secret")

func TestRegister_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, jwtSecret)

	email := "test@example.com"
	password := "securepass"
	role := "employee"

	repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		return user.Email == email && user.Role == role && err == nil
	})).Return(nil)

	token, err := svc.Register(context.Background(), email, password, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, role, claims["user_type"])

	repo.AssertExpectations(t)
}

func TestRegister_InvalidInput(t *testing.T) {
	svc := service.NewAuthService(nil, jwtSecret)

	token, err := svc.Register(context.Background(), "", "pass", "employee")
	assert.Empty(t, token)
	assert.EqualError(t, err, "email/password required")
}

func TestRegister_CreateUserError(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, jwtSecret)

	email := "test@example.com"
	password := "securepass"
	role := "employee"

	repo.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("db error"))

	token, err := svc.Register(context.Background(), email, password, role)
	assert.Empty(t, token)
	assert.EqualError(t, err, "db error")
}

func TestLogin_Success(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, jwtSecret)

	email := "test@example.com"
	password := "mypassword"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &domain.User{
		Email:    email,
		Password: string(hashed),
		Role:     "employee",
	}

	repo.On("GetByEmail", mock.Anything, email).Return(user, nil)

	token, err := svc.Login(context.Background(), email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	repo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, jwtSecret)

	repo.On("GetByEmail", mock.Anything, "no@user.com").Return(nil, errors.New("not found"))

	token, err := svc.Login(context.Background(), "no@user.com", "pass")
	assert.Empty(t, token)
	assert.EqualError(t, err, "not found")
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := new(mockUserRepo)
	svc := service.NewAuthService(repo, jwtSecret)

	email := "test@example.com"
	wrongPassword := "wrongpass"
	correctPassword := "correct"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	user := &domain.User{Email: email, Password: string(hashed), Role: "employee"}

	repo.On("GetByEmail", mock.Anything, email).Return(user, nil)

	token, err := svc.Login(context.Background(), email, wrongPassword)
	assert.Empty(t, token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hashedPassword")

	repo.AssertExpectations(t)
}
