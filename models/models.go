package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type BaseModel struct {
	Id uuid.UUID `json:"id" dynamo:"id"`

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}
