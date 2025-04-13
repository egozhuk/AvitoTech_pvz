package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"pvs/internal/domain"
	"pvs/internal/repository"
)

type ProductService struct {
	productRepo   repository.ProductRepository
	receptionRepo repository.ReceptionRepository
}

func NewProductService(productRepo repository.ProductRepository, receptionRepo repository.ReceptionRepository) *ProductService {
	return &ProductService{productRepo: productRepo, receptionRepo: receptionRepo}
}

func (s *ProductService) AddProduct(ctx context.Context, pvzID uuid.UUID, productType string, role string) (*domain.Product, error) {
	if role != "employee" {
		return nil, errors.New("только сотрудники могут добавлять товары")
	}

	reception, err := s.receptionRepo.GetOpenReception(ctx, pvzID)
	if err != nil {
		return nil, errors.New("нет активной приёмки")
	}

	return s.productRepo.AddProduct(ctx, reception.ID, productType)
}

func (s *ProductService) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID, role string) error {
	if role != "employee" {
		return errors.New("только сотрудники могут удалять товары")
	}

	reception, err := s.receptionRepo.GetOpenReception(ctx, pvzID)
	if err != nil {
		return errors.New("нет активной приёмки")
	}

	return s.productRepo.DeleteLastProduct(ctx, reception.ID)
}
