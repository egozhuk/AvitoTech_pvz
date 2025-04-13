package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Type        string
	ReceptionID uuid.UUID
	DateTime    time.Time
}
