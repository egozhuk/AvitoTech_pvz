package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"pvs/internal/domain"
)

type PostgresReceptionRepository struct {
	pool *pgxpool.Pool
}

func NewReceptionRepository(pool *pgxpool.Pool) *PostgresReceptionRepository {
	return &PostgresReceptionRepository{pool: pool}
}

func (r *PostgresReceptionRepository) CreateReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
	var rec domain.Reception
	err := r.pool.QueryRow(ctx,
		`INSERT INTO reception (id, pvz_id, status) VALUES (gen_random_uuid(), $1, 'in_progress')
		 RETURNING id, pvz_id, date_time, status`,
		pvzID,
	).Scan(&rec.ID, &rec.PVZID, &rec.DateTime, &rec.Status)
	return &rec, err
}

func (r *PostgresReceptionRepository) GetOpenReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
	var rec domain.Reception
	err := r.pool.QueryRow(ctx, `
		SELECT id, pvz_id, date_time, status FROM reception
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC LIMIT 1
	`, pvzID).Scan(&rec.ID, &rec.PVZID, &rec.DateTime, &rec.Status)
	return &rec, err
}

func (r *PostgresReceptionRepository) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
	var rec domain.Reception
	err := r.pool.QueryRow(ctx, `
		UPDATE reception
		SET status = 'close'
		WHERE id = (
			SELECT id FROM reception
			WHERE pvz_id = $1 AND status = 'in_progress'
			ORDER BY date_time DESC
			LIMIT 1
		)
		RETURNING id, pvz_id, date_time, status
	`, pvzID).Scan(&rec.ID, &rec.PVZID, &rec.DateTime, &rec.Status)
	return &rec, err
}
