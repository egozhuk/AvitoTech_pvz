package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"pvs/internal/domain"
)

type PostgresProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{pool: pool}
}

func (r *PostgresProductRepository) AddProduct(ctx context.Context, receptionID uuid.UUID, productType string) (*domain.Product, error) {
	var p domain.Product
	err := r.pool.QueryRow(ctx, `
		INSERT INTO product (id, type, reception_id)
		VALUES (gen_random_uuid(), $1, $2)
		RETURNING id, type, reception_id, date_time
	`, productType, receptionID).Scan(&p.ID, &p.Type, &p.ReceptionID, &p.DateTime)
	return &p, err
}

func (r *PostgresProductRepository) DeleteLastProduct(ctx context.Context, receptionID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM product
		WHERE id = (
			SELECT id FROM product
			WHERE reception_id = $1
			ORDER BY date_time DESC LIMIT 1
		)
	`, receptionID)
	return err
}

func (r *PostgresProductRepository) GetProductsByReception(ctx context.Context, receptionID uuid.UUID) ([]domain.Product, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, type, reception_id, date_time FROM product
		WHERE reception_id = $1
		ORDER BY date_time ASC
	`, receptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Type, &p.ReceptionID, &p.DateTime); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
