package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvs/internal/controller"
	"pvs/internal/domain"
)

type mockReceptionService struct {
	mock.Mock
}

func (m *mockReceptionService) CreateReception(ctx context.Context, pvzID uuid.UUID, role string) (*domain.Reception, error) {
	args := m.Called(ctx, pvzID, role)
	if rec := args.Get(0); rec != nil {
		return rec.(*domain.Reception), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockReceptionService) CloseLastReception(ctx context.Context, pvzID uuid.UUID, role string) (*domain.Reception, error) {
	args := m.Called(ctx, pvzID, role)
	if rec := args.Get(0); rec != nil {
		return rec.(*domain.Reception), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestCreateReceptionHandler_BadJSON(t *testing.T) {
	handler := controller.CreateReceptionHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/reception", bytes.NewReader([]byte("bad json")))
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCloseLastReceptionHandler_Success(t *testing.T) {
	service := new(mockReceptionService)
	handler := controller.CloseLastReceptionHandler(service)

	pvzID := uuid.New()
	expected := &domain.Reception{ID: uuid.New(), PVZID: pvzID}
	service.On("CloseLastReception", mock.Anything, pvzID, "employee").Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/reception/"+pvzID.String()+"/close", nil)
	req.SetPathValue("pvzId", pvzID.String())
	req = req.WithContext(withRole(req.Context(), "employee"))
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp domain.Reception
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, expected.ID, resp.ID)
}

func TestCloseLastReceptionHandler_InvalidUUID(t *testing.T) {
	handler := controller.CloseLastReceptionHandler(nil)

	req := httptest.NewRequest(http.MethodPost, "/reception/invalid/close", nil)
	req.SetPathValue("pvzId", "not-a-uuid")
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCloseLastReceptionHandler_ServiceError(t *testing.T) {
	service := new(mockReceptionService)
	handler := controller.CloseLastReceptionHandler(service)

	pvzID := uuid.New()
	service.On("CloseLastReception", mock.Anything, pvzID, "employee").Return(nil, errors.New("fail"))

	req := httptest.NewRequest(http.MethodPost, "/reception/"+pvzID.String()+"/close", nil)
	req.SetPathValue("pvzId", pvzID.String())
	req = req.WithContext(withRole(req.Context(), "employee"))
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
