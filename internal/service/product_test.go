package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvs/internal/domain"
	"pvs/internal/service"
)

type mockProductRepo struct {
	mock.Mock
}

func (m *mockProductRepo) AddProduct(ctx context.Context, receptionID uuid.UUID, productType string) (*domain.Product, error) {
	args := m.Called(ctx, receptionID, productType)
	if prod := args.Get(0); prod != nil {
		return prod.(*domain.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProductRepo) DeleteLastProduct(ctx context.Context, receptionID uuid.UUID) error {
	args := m.Called(ctx, receptionID)
	return args.Error(0)
}

func (m *mockProductRepo) GetProductsByReception(ctx context.Context, receptionID uuid.UUID) ([]domain.Product, error) {
	args := m.Called(ctx, receptionID)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func TestAddProduct_Success(t *testing.T) {
	productRepo := new(mockProductRepo)
	receptionRepo := new(mockReceptionRepo)
	svc := service.NewProductService(productRepo, receptionRepo)

	pvzID := uuid.New()
	receptionID := uuid.New()
	productType := "electronics"
	expectedProduct := &domain.Product{ID: uuid.New(), Type: productType}

	receptionRepo.On("GetOpenReception", mock.Anything, pvzID).Return(&domain.Reception{ID: receptionID}, nil)
	productRepo.On("AddProduct", mock.Anything, receptionID, productType).Return(expectedProduct, nil)

	result, err := svc.AddProduct(context.Background(), pvzID, productType, "employee")
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, result)

	receptionRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestAddProduct_Unauthorized(t *testing.T) {
	svc := service.NewProductService(nil, nil)

	result, err := svc.AddProduct(context.Background(), uuid.New(), "toys", "moderator")
	assert.Nil(t, result)
	assert.EqualError(t, err, "только сотрудники могут добавлять товары")
}

func TestAddProduct_NoReception(t *testing.T) {
	productRepo := new(mockProductRepo)
	receptionRepo := new(mockReceptionRepo)
	svc := service.NewProductService(productRepo, receptionRepo)

	pvzID := uuid.New()

	receptionRepo.On("GetOpenReception", mock.Anything, pvzID).Return(nil, errors.New("not found"))

	result, err := svc.AddProduct(context.Background(), pvzID, "books", "employee")
	assert.Nil(t, result)
	assert.EqualError(t, err, "нет активной приёмки")

	receptionRepo.AssertExpectations(t)
}

func TestDeleteLastProduct_Success(t *testing.T) {
	productRepo := new(mockProductRepo)
	receptionRepo := new(mockReceptionRepo)
	svc := service.NewProductService(productRepo, receptionRepo)

	pvzID := uuid.New()
	receptionID := uuid.New()

	receptionRepo.On("GetOpenReception", mock.Anything, pvzID).Return(&domain.Reception{ID: receptionID}, nil)
	productRepo.On("DeleteLastProduct", mock.Anything, receptionID).Return(nil)

	err := svc.DeleteLastProduct(context.Background(), pvzID, "employee")
	assert.NoError(t, err)

	receptionRepo.AssertExpectations(t)
	productRepo.AssertExpectations(t)
}

func TestDeleteLastProduct_Unauthorized(t *testing.T) {
	svc := service.NewProductService(nil, nil)

	err := svc.DeleteLastProduct(context.Background(), uuid.New(), "moderator")
	assert.EqualError(t, err, "только сотрудники могут удалять товары")
}

func TestDeleteLastProduct_NoReception(t *testing.T) {
	productRepo := new(mockProductRepo)
	receptionRepo := new(mockReceptionRepo)
	svc := service.NewProductService(productRepo, receptionRepo)

	pvzID := uuid.New()

	receptionRepo.On("GetOpenReception", mock.Anything, pvzID).Return(nil, errors.New("none"))

	err := svc.DeleteLastProduct(context.Background(), pvzID, "employee")
	assert.EqualError(t, err, "нет активной приёмки")

	receptionRepo.AssertExpectations(t)
}
