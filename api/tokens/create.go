package tokens

import (
	"context"
)

type Handler struct {
}

type NewToken struct {
}

type NewTokenRequest struct {
	Body NewToken `json:"token"`
}

type NewTokenResponse struct{}

func (h *Handler) CreateToken(ctx context.Context, in *NewTokenRequest) (*NewTokenResponse, error) {

}
