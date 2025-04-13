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

type mockReceptionRepo struct {
	mock.Mock
}

func (m *mockReceptionRepo) GetOpenReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
	args := m.Called(ctx, pvzID)
	if rec := args.Get(0); rec != nil {
		return rec.(*domain.Reception), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockReceptionRepo) CreateReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(*domain.Reception), args.Error(1)
}

func (m *mockReceptionRepo) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(*domain.Reception), args.Error(1)
}

func TestCreateReception_Success(t *testing.T) {
	repo := new(mockReceptionRepo)
	svc := service.NewReceptionService(repo)

	pvzID := uuid.New()
	expected := &domain.Reception{ID: uuid.New(), PVZID: pvzID}

	repo.On("GetOpenReception", mock.Anything, pvzID).Return(nil, errors.New("not found"))
	repo.On("CreateReception", mock.Anything, pvzID).Return(expected, nil)

	rec, err := svc.CreateReception(context.Background(), pvzID, "employee")
	assert.NoError(t, err)
	assert.Equal(t, expected, rec)
	repo.AssertExpectations(t)
}

func TestCreateReception_AlreadyOpen(t *testing.T) {
	repo := new(mockReceptionRepo)
	svc := service.NewReceptionService(repo)

	pvzID := uuid.New()
	existing := &domain.Reception{ID: uuid.New(), PVZID: pvzID}

	repo.On("GetOpenReception", mock.Anything, pvzID).Return(existing, nil)

	rec, err := svc.CreateReception(context.Background(), pvzID, "employee")
	assert.Nil(t, rec)
	assert.EqualError(t, err, "уже есть незакрытая приёмка")
	repo.AssertExpectations(t)
}

func TestCreateReception_NotEmployee(t *testing.T) {
	repo := new(mockReceptionRepo)
	svc := service.NewReceptionService(repo)

	pvzID := uuid.New()
	rec, err := svc.CreateReception(context.Background(), pvzID, "moderator")
	assert.Nil(t, rec)
	assert.EqualError(t, err, "доступ разрешён только сотрудникам ПВЗ")
}

func TestCloseLastReception_Success(t *testing.T) {
	repo := new(mockReceptionRepo)
	svc := service.NewReceptionService(repo)

	pvzID := uuid.New()
	expected := &domain.Reception{ID: uuid.New(), PVZID: pvzID}

	repo.On("CloseLastReception", mock.Anything, pvzID).Return(expected, nil)

	rec, err := svc.CloseLastReception(context.Background(), pvzID, "employee")
	assert.NoError(t, err)
	assert.Equal(t, expected, rec)
	repo.AssertExpectations(t)
}

func TestCloseLastReception_NotEmployee(t *testing.T) {
	repo := new(mockReceptionRepo)
	svc := service.NewReceptionService(repo)

	pvzID := uuid.New()
	rec, err := svc.CloseLastReception(context.Background(), pvzID, "moderator")
	assert.Nil(t, rec)
	assert.EqualError(t, err, "доступ разрешён только сотрудникам ПВЗ")
}
