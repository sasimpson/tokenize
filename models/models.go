package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	Id uuid.UUID `json:"id"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
