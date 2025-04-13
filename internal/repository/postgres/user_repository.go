package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"pvs/internal/domain"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, role) VALUES (gen_random_uuid(), $1, $2, $3)`,
		u.Email, u.Password, u.Role)
	return err
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, email, password_hash, role FROM users WHERE email=$1`, email)
	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	return &user, err
}
