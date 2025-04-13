package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"pvs/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type PVZRepository interface {
	CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error)
	ListPVZWithFilter(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]domain.PVZ, error)
}

type ReceptionRepository interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error)
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error)
	GetOpenReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error)
}

type ProductRepository interface {
	AddProduct(ctx context.Context, receptionID uuid.UUID, productType string) (*domain.Product, error)
	DeleteLastProduct(ctx context.Context, receptionID uuid.UUID) error
	GetProductsByReception(ctx context.Context, receptionID uuid.UUID) ([]domain.Product, error)
}
