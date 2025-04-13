package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pvs/internal/domain"
)

type PostgresPVZRepository struct {
	pool *pgxpool.Pool
}

func NewPVSRepository(pool *pgxpool.Pool) *PostgresPVZRepository {
	return &PostgresPVZRepository{pool: pool}
}

func (r *PostgresPVZRepository) CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO pvz (id, city) VALUES (gen_random_uuid(), $1)
		RETURNING id, city, registration_date
	`, city)

	var pvz domain.PVZ
	err := row.Scan(&pvz.ID, &pvz.City, &pvz.RegistrationDate)
	return &pvz, err
}

func (r *PostgresPVZRepository) ListPVZWithFilter(
	ctx context.Context,
	startDate, endDate *time.Time,
	page, limit int,
) ([]domain.PVZ, error) {
	query := `
		SELECT DISTINCT pvz.id, pvz.city, pvz.registration_date
		FROM pvz
		LEFT JOIN reception r ON r.pvz_id = pvz.id
		WHERE ($1::timestamp IS NULL OR r.date_time >= $1::timestamp)
		  AND ($2::timestamp IS NULL OR r.date_time <= $2::timestamp)
		ORDER BY pvz.registration_date DESC
		LIMIT $3 OFFSET $4
	`

	offset := (page - 1) * limit

	log.Printf("📦 ListPVZWithFilter called with: startDate=%v, endDate=%v, page=%d, limit=%d", startDate, endDate, page, limit)

	rows, err := r.pool.Query(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		log.Printf("❌ query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var result []domain.PVZ
	for rows.Next() {
		var p domain.PVZ
		if err := rows.Scan(&p.ID, &p.City, &p.RegistrationDate); err != nil {
			log.Printf("❌ scan error: %v", err)
			return nil, err
		}
		result = append(result, p)
	}

	log.Printf("✅ ListPVZWithFilter result count: %d", len(result))
	return result, nil
}
