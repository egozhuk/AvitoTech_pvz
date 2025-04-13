package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"pvs/internal/domain"
	"pvs/internal/transport/middleware"
	"strconv"
	"time"
)

type CreatePVZRequest struct {
	City string `json:"city"`
}

type PVZServiceInterface interface {
	CreatePVZ(ctx context.Context, city string) (*domain.PVZ, error)
	ListPVZWithFilter(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]domain.PVZ, error)
}

func CreatePVZHandler(s PVZServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := middleware.GetUserRole(r.Context())
		if role != "moderator" {
			http.Error(w, "только модератор может создавать ПВЗ", http.StatusForbidden)
			return
		}

		var req CreatePVZRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "неверный формат запроса", http.StatusBadRequest)
			return
		}

		pvz, err := s.CreatePVZ(r.Context(), req.City)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pvz)
	}
}

func GetPVZListHandler(s PVZServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		var from, to *time.Time

		if v := query.Get("startDate"); v != "" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				http.Error(w, "неверный формат startDate", http.StatusBadRequest)
				return
			}
			from = &t
		}
		if v := query.Get("endDate"); v != "" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				http.Error(w, "неверный формат endDate", http.StatusBadRequest)
				return
			}
			to = &t
		}

		page, _ := strconv.Atoi(query.Get("page"))
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(query.Get("limit"))
		if limit < 1 || limit > 30 {
			limit = 10
		}

		result, err := s.ListPVZWithFilter(r.Context(), from, to, page, limit)
		if err != nil {
			http.Error(w, "ошибка получения списка ПВЗ", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}
