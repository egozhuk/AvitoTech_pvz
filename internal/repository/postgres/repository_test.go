package postgres_test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"pvs/internal/domain"
	"pvs/internal/repository/postgres"
	"testing"
	"time"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "pvs_db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	defer container.Terminate(ctx)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432/tcp")
	dsn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/pvs_db?sslmode=disable", host, port.Port())

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";
		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL
		);
		CREATE TABLE pvz (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			city TEXT NOT NULL,
			registration_date TIMESTAMP NOT NULL DEFAULT now()
		);
		CREATE TABLE reception (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			pvz_id UUID NOT NULL REFERENCES pvz(id),
			date_time TIMESTAMP NOT NULL DEFAULT now(),
			status TEXT NOT NULL DEFAULT 'in_progress'
		);
		CREATE TABLE product (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			type TEXT NOT NULL,
			reception_id UUID NOT NULL REFERENCES reception(id),
			date_time TIMESTAMP NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		panic(err)
	}

	testDB = pool
	code := m.Run()
	pool.Close()
	os.Exit(code)
}

func TestAllRepositories(t *testing.T) {
	ctx := context.Background()

	userRepo := postgres.NewUserRepository(testDB)
	pvzRepo := postgres.NewPVSRepository(testDB)
	receptionRepo := postgres.NewReceptionRepository(testDB)
	productRepo := postgres.NewProductRepository(testDB)

	// Create user
	err := userRepo.CreateUser(ctx, &domain.User{
		Email:    "test@mail.com",
		Password: "hashed",
		Role:     "employee",
	})
	require.NoError(t, err)

	user, err := userRepo.GetByEmail(ctx, "test@mail.com")
	require.NoError(t, err)
	assert.Equal(t, "test@mail.com", user.Email)

	// Create PVZ
	pvz, err := pvzRepo.CreatePVZ(ctx, "Москва")
	require.NoError(t, err)

	// Create Reception
	reception, err := receptionRepo.CreateReception(ctx, pvz.ID)
	require.NoError(t, err)
	assert.Equal(t, "in_progress", reception.Status)

	// Add product
	product, err := productRepo.AddProduct(ctx, reception.ID, "book")
	require.NoError(t, err)
	assert.Equal(t, "book", product.Type)

	// Get products
	products, err := productRepo.GetProductsByReception(ctx, reception.ID)
	require.NoError(t, err)
	assert.Len(t, products, 1)

	// Delete last product
	err = productRepo.DeleteLastProduct(ctx, reception.ID)
	require.NoError(t, err)

	// Close reception
	closed, err := receptionRepo.CloseLastReception(ctx, pvz.ID)
	require.NoError(t, err)
	assert.Equal(t, "close", closed.Status)

	// List PVZ
	pvzList, err := pvzRepo.ListPVZWithFilter(ctx, nil, nil, 1, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, pvzList)
}
