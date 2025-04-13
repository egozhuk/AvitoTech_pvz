package controller_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"pvs/internal/controller"
	"pvs/internal/domain"
	"pvs/internal/transport/middleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPVZService struct {
	mock.Mock
}

func (m *mockPVZService) CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error) {
	args := m.Called(ctx, city)
	if pvz := args.Get(0); pvz != nil {
		return pvz.(*domain.PVZ), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPVZService) ListPVZWithFilter(ctx context.Context, from, to *time.Time, page, limit int) ([]domain.PVZ, error) {
	args := m.Called(ctx, from, to, page, limit)
	return args.Get(0).([]domain.PVZ), args.Error(1)
}

func withRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, middleware.RoleCtxKey{}, role)
}

func TestCreatePVZHandler_Forbidden(t *testing.T) {
	handler := controller.CreatePVZHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader([]byte(`{"city":"Москва"}`)))
	req = req.WithContext(withRole(req.Context(), "employee")) // not moderator
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreatePVZHandler_InvalidJSON(t *testing.T) {
	handler := controller.CreatePVZHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader([]byte(`not-json`)))
	req = req.WithContext(withRole(req.Context(), "moderator"))
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePVZHandler_ServiceError(t *testing.T) {
	service := new(mockPVZService)
	handler := controller.CreatePVZHandler(service)

	service.On("CreatePVZ", mock.Anything, "Казань").Return(nil, errors.New("city error"))

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader([]byte(`{"city":"Казань"}`)))
	req = req.WithContext(withRole(req.Context(), "moderator"))
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPVZListHandler_InvalidStartDate(t *testing.T) {
	handler := controller.GetPVZListHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/pvz?startDate=invalid", nil)
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
