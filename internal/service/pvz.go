package service

import (
	"context"
	"errors"
	"time"

	"pvs/internal/domain"
	"pvs/internal/repository"
)

type PVZService struct {
	repo repository.PVZRepository
}

func NewPVSService(repo repository.PVZRepository) *PVZService {
	return &PVZService{repo: repo}
}

func (s *PVZService) CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error) {
	allowed := map[string]struct{}{"Москва": {}, "Санкт-Петербург": {}, "Казань": {}}
	if _, ok := allowed[city]; !ok {
		return nil, errors.New("город недоступен для регистрации")
	}
	return s.repo.CreatePVZ(ctx, city)
}

func (s *PVZService) ListPVZWithFilter(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]domain.PVZ, error) {
	return s.repo.ListPVZWithFilter(ctx, startDate, endDate, page, limit)
}
