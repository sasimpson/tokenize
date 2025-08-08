package api

import (
	"context"

	"tokenize/api"
)

type TokenHandler struct {
	api.BaseHandler
}

type NewToken struct {
}

type NewTokenRequest struct {
	Body NewToken `json:"token"`
}

type NewTokenResponse struct{}

func (h *TokenHandler) CreateToken(ctx context.Context, in *NewTokenRequest) (*NewTokenResponse, error) {

}
