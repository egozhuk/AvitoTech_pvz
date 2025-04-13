package controller

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"pvs/internal/domain"
	"pvs/internal/transport/middleware"
)

type ReceptionCreateRequest struct {
	PVZID uuid.UUID `json:"pvzId"`
}

type ReceptionServiceInterface interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID, role string) (*domain.Reception, error)
	CloseLastReception(ctx context.Context, pvzID uuid.UUID, role string) (*domain.Reception, error)
}

func CreateReceptionHandler(s ReceptionServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ReceptionCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "неверный формат запроса", http.StatusBadRequest)
			return
		}

		role := middleware.GetUserRole(r.Context())
		reception, err := s.CreateReception(r.Context(), req.PVZID, role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(reception)
	}
}

func CloseLastReceptionHandler(s ReceptionServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pvzIDStr := r.PathValue("pvzId")
		pvzID, err := uuid.Parse(pvzIDStr)
		if err != nil {
			http.Error(w, "неверный UUID", http.StatusBadRequest)
			return
		}
		role := middleware.GetUserRole(r.Context())
		reception, err := s.CloseLastReception(r.Context(), pvzID, role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(reception)
	}
}
