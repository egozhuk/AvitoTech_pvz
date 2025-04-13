package controller

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"pvs/internal/domain"
	"pvs/internal/transport/middleware"
)

type AddProductRequest struct {
	Type  string    `json:"type"`
	PVZID uuid.UUID `json:"pvzId"`
}

type ProductServiceInterface interface {
	AddProduct(ctx context.Context, pvzID uuid.UUID, productType, role string) (*domain.Product, error)
	DeleteLastProduct(ctx context.Context, pvzID uuid.UUID, role string) error
}

func AddProductHandler(s ProductServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "неверный запрос", http.StatusBadRequest)
			return
		}
		role := middleware.GetUserRole(r.Context())
		product, err := s.AddProduct(r.Context(), req.PVZID, req.Type, role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

func DeleteLastProductHandler(s ProductServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pvzIDStr := r.PathValue("pvzId")
		pvzID, err := uuid.Parse(pvzIDStr)
		if err != nil {
			http.Error(w, "неверный UUID", http.StatusBadRequest)
			return
		}
		role := middleware.GetUserRole(r.Context())
		if err := s.DeleteLastProduct(r.Context(), pvzID, role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
