package persistence

import (
	"context"

	"tokenize/models"
)

type Store interface {
	GetToken(context.Context, string) (*models.Token, error)
	CreateToken(context.Context, *models.Token) (*models.Token, error)
	DeleteToken(context.Context, *models.Token) error
}
