package domain

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID
	PVZID    uuid.UUID
	DateTime time.Time
	Status   string
}
