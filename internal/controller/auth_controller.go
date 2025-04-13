package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"pvs/internal/service"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, email, password, role string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

type AuthTokenResponse struct {
	Token string `json:"token"`
}

func DummyLoginHandler(secret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DummyLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if req.Role != "employee" && req.Role != "moderator" {
			http.Error(w, "invalid role", http.StatusBadRequest)
			return
		}

		token, err := service.GenerateToken(secret, req.Role)
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(AuthTokenResponse{Token: token})
	}
}

func RegisterHandler(auth AuthServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if req.Role != "employee" && req.Role != "moderator" {
			http.Error(w, "invalid role", http.StatusBadRequest)
			return
		}
		token, err := auth.Register(r.Context(), req.Email, req.Password, req.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(AuthTokenResponse{Token: token})
	}
}

func LoginHandler(auth AuthServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		token, err := auth.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(AuthTokenResponse{Token: token})
	}
}
