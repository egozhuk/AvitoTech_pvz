package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"pvs/internal/config"
	"pvs/internal/controller"
	"pvs/internal/repository/postgres"
	"pvs/internal/service"
	"pvs/internal/transport/middleware"
)

type App struct {
	Server *http.Server
	DB     *pgxpool.Pool
}

func NewApp(ctx context.Context) (*App, error) {
	cfg := config.Load()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	execPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd failed: %w", err)
	}

	migrationsPath := filepath.Join(execPath, "migrations")

	sqlDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	defer sqlDB.Close()

	log.Println("📦 applying migrations from:", migrationsPath)

	if err := goose.Up(sqlDB, migrationsPath); err != nil {
		return nil, fmt.Errorf("goose.Up: %w", err)
	}

	userRepo := postgres.NewUserRepository(db)
	pvzRepo := postgres.NewPVSRepository(db)
	receptionRepo := postgres.NewReceptionRepository(db)
	productRepo := postgres.NewProductRepository(db)

	authService := service.NewAuthService(userRepo, []byte(cfg.JWTSecret))
	pvzService := service.NewPVSService(pvzRepo)
	receptionService := service.NewReceptionService(receptionRepo)
	productService := service.NewProductService(productRepo, receptionRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /dummyLogin", controller.DummyLoginHandler([]byte(cfg.JWTSecret)))
	mux.HandleFunc("POST /register", controller.RegisterHandler(authService))
	mux.HandleFunc("POST /login", controller.LoginHandler(authService))

	protected := http.NewServeMux()
	protected.HandleFunc("POST /pvz", controller.CreatePVZHandler(pvzService))
	protected.HandleFunc("GET /pvz", controller.GetPVZListHandler(pvzService))

	protected.HandleFunc("POST /receptions", controller.CreateReceptionHandler(receptionService))
	protected.HandleFunc("POST /pvz/{pvzId}/close_last_reception", controller.CloseLastReceptionHandler(receptionService))

	protected.HandleFunc("POST /products", controller.AddProductHandler(productService))
	protected.HandleFunc("POST /pvz/{pvzId}/delete_last_product", controller.DeleteLastProductHandler(productService))

	mux.Handle("/", middleware.AuthMiddleware(protected))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &App{Server: server, DB: db}, nil
}

func (a *App) Run() error {
	log.Println("✅ Server running on :8080")
	return a.Server.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	log.Println("🛑 Shutting down...")
	a.DB.Close()
	return a.Server.Shutdown(ctx)
}
