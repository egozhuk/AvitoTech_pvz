package domain

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID               uuid.UUID
	City             string
	RegistrationDate time.Time
}
