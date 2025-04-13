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

type mockProductService struct {
	mock.Mock
}

func (m *mockProductService) AddProduct(ctx context.Context, pvzID uuid.UUID, productType, role string) (*domain.Product, error) {
	args := m.Called(ctx, pvzID, productType, role)
	if p := args.Get(0); p != nil {
		return p.(*domain.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProductService) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID, role string) error {
	args := m.Called(ctx, pvzID, role)
	return args.Error(0)
}

func TestAddProductHandler_BadRequest(t *testing.T) {
	handler := controller.AddProductHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddProductHandler_Success(t *testing.T) {
	service := new(mockProductService)
	handler := controller.AddProductHandler(service)

	pvzID := uuid.New()
	expected := &domain.Product{
		ID:          uuid.New(),
		Type:        "shoes",
		ReceptionID: uuid.New(),
	}
	service.On("AddProduct", mock.Anything, pvzID, "shoes", "employee").Return(expected, nil)

	body, _ := json.Marshal(map[string]any{
		"pvzId": pvzID,
		"type":  "shoes",
	})

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req = req.WithContext(withRole(req.Context(), "employee"))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var result domain.Product
	err := json.NewDecoder(w.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestAddProductHandler_ServiceError(t *testing.T) {
	service := new(mockProductService)
	handler := controller.AddProductHandler(service)

	pvzID := uuid.New()
	service.On("AddProduct", mock.Anything, pvzID, "toys", "employee").Return(nil, errors.New("fail"))

	body, _ := json.Marshal(map[string]any{
		"pvzId": pvzID,
		"type":  "toys",
	})

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req = req.WithContext(withRole(req.Context(), "employee"))
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "fail")
}

func TestDeleteLastProductHandler_BadUUID(t *testing.T) {
	handler := controller.DeleteLastProductHandler(nil)
	req := httptest.NewRequest(http.MethodPost, "/pvz/invalid-uuid/delete_last_product", nil)
	req.SetPathValue("pvzId", "invalid-uuid")
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteLastProductHandler_Success(t *testing.T) {
	service := new(mockProductService)
	handler := controller.DeleteLastProductHandler(service)

	pvzID := uuid.New()
	service.On("DeleteLastProduct", mock.Anything, pvzID, "employee").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)
	req.SetPathValue("pvzId", pvzID.String())
	req = req.WithContext(withRole(req.Context(), "employee"))
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteLastProductHandler_ServiceError(t *testing.T) {
	service := new(mockProductService)
	handler := controller.DeleteLastProductHandler(service)

	pvzID := uuid.New()
	service.On("DeleteLastProduct", mock.Anything, pvzID, "employee").Return(errors.New("fail"))

	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)
	req.SetPathValue("pvzId", pvzID.String())
	req = req.WithContext(withRole(req.Context(), "employee"))
	w := httptest.NewRecorder()

	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "fail")
}
