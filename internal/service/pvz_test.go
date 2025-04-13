package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pvs/internal/domain"
	"pvs/internal/service"
)

type mockPVZRepo struct {
	mock.Mock
}

func (m *mockPVZRepo) CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error) {
	args := m.Called(ctx, city)
	if pvz := args.Get(0); pvz != nil {
		return pvz.(*domain.PVZ), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPVZRepo) ListPVZWithFilter(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]domain.PVZ, error) {
	args := m.Called(ctx, startDate, endDate, page, limit)
	return args.Get(0).([]domain.PVZ), args.Error(1)
}

func TestCreatePVZ_AllowedCity(t *testing.T) {
	repo := new(mockPVZRepo)
	svc := service.NewPVSService(repo)

	city := "Москва"
	expected := &domain.PVZ{City: city}

	repo.On("CreatePVZ", mock.Anything, city).Return(expected, nil)

	pvz, err := svc.CreatePVZ(context.Background(), city)
	assert.NoError(t, err)
	assert.Equal(t, expected, pvz)
	repo.AssertExpectations(t)
}

func TestCreatePVZ_DisallowedCity(t *testing.T) {
	repo := new(mockPVZRepo)
	svc := service.NewPVSService(repo)

	pvz, err := svc.CreatePVZ(context.Background(), "Новосибирск")
	assert.Nil(t, pvz)
	assert.EqualError(t, err, "город недоступен для регистрации")
}

func TestListPVZWithFilter(t *testing.T) {
	repo := new(mockPVZRepo)
	svc := service.NewPVSService(repo)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	page := 1
	limit := 10

	expected := []domain.PVZ{
		{City: "Москва"},
		{City: "Казань"},
	}

	repo.On("ListPVZWithFilter", mock.Anything, &start, &end, page, limit).Return(expected, nil)

	result, err := svc.ListPVZWithFilter(context.Background(), &start, &end, page, limit)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}
