package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	Id uuid.UUID `json:"id" dynamo:"id"`

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}
