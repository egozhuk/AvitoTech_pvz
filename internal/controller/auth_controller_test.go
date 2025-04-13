package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pvs/internal/controller"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, email, password, role string) (string, error) {
	args := m.Called(ctx, email, password, role)
	return args.String(0), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestDummyLoginHandler_Success(t *testing.T) {
	body := []byte(`{"role":"employee"}`)
	req := httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := controller.DummyLoginHandler([]byte("secret"))
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp controller.AuthTokenResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.NotEmpty(t, resp.Token)
}

func TestRegisterHandler_Success(t *testing.T) {
	auth := new(mockAuthService)
	auth.On("Register", mock.Anything, "test@mail.com", "123", "employee").Return("token123", nil)

	reqBody := []byte(`{"email":"test@mail.com","password":"123","role":"employee"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler := controller.RegisterHandler(auth)
	handler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp controller.AuthTokenResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "token123", resp.Token)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	auth := new(mockAuthService)
	auth.On("Login", mock.Anything, "test@mail.com", "wrongpass").Return("", errors.New("invalid"))

	reqBody := []byte(`{"email":"test@mail.com","password":"wrongpass"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler := controller.LoginHandler(auth)
	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
