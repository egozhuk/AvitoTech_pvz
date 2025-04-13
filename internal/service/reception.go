package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"pvs/internal/domain"
	"pvs/internal/repository"
)

type ReceptionService struct {
	repo repository.ReceptionRepository
}

func NewReceptionService(repo repository.ReceptionRepository) *ReceptionService {
	return &ReceptionService{repo: repo}
}

func (s *ReceptionService) CreateReception(ctx context.Context, pvzID uuid.UUID, role string) (*domain.Reception, error) {
	if role != "employee" {
		return nil, errors.New("доступ разрешён только сотрудникам ПВЗ")
	}
	open, err := s.repo.GetOpenReception(ctx, pvzID)
	if err == nil && open != nil {
		return nil, errors.New("уже есть незакрытая приёмка")
	}
	return s.repo.CreateReception(ctx, pvzID)
}

func (s *ReceptionService) CloseLastReception(ctx context.Context, pvzID uuid.UUID, role string) (*domain.Reception, error) {
	if role != "employee" {
		return nil, errors.New("доступ разрешён только сотрудникам ПВЗ")
	}
	return s.repo.CloseLastReception(ctx, pvzID)
}
