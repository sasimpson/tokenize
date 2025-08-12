package mock

import (
	"context"

	"tokenize/models"
)

type Store struct {
	Token       *models.Token
	CreateError error
	GetError    error
	DeleteError error
}

func (s Store) GetToken(_ context.Context, _ string) (*models.Token, error) {
	return s.Token, s.GetError
}

func (s Store) CreateToken(_ context.Context, _ *models.Token) (*models.Token, error) {
	return s.Token, s.CreateError
}

func (s Store) DeleteToken(_ context.Context, _ *models.Token) error {
	return s.DeleteError
}
