package api

import (
	"context"
	"net/http"
	"time"
	"tokenize/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

func (h *BaseHandler) RegisterTokensRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "CreateToken",
		Summary:       "Create a new token",
		Method:        http.MethodPost,
		Path:          "/token",
		DefaultStatus: http.StatusCreated,
		Errors: []int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusBadRequest,
		},
	}, h.CreateToken)

	huma.Register(api, huma.Operation{
		OperationID:   "GetEncryptedToken",
		Summary:       "Get an encrypted token properties",
		Method:        http.MethodGet,
		Path:          "/token/{token}",
		DefaultStatus: http.StatusOK,
		Errors: []int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusBadRequest,
			http.StatusNotFound,
		},
	}, h.GetEncryptedToken)

	huma.Register(api, huma.Operation{
		OperationID:   "GetDecryptedToken",
		Summary:       "Get a decrypted token and properties",
		Method:        http.MethodGet,
		Path:          "/token/{token}/decrypt",
		DefaultStatus: http.StatusOK,
		Errors: []int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusBadRequest,
			http.StatusNotFound,
		},
	}, h.GetDecryptedToken)

}

type NewTokenRequest struct {
	Body struct {
		Data models.CreateToken `json:"data" validate:"required"`
	}
}

type NewTokenResponse struct {
	Body struct {
		Token string `json:"token"`
	}
}

func (h *BaseHandler) CreateToken(_ context.Context, in *NewTokenRequest) (*NewTokenResponse, error) {

	newToken := models.Token{
		CreateToken: in.Body.Data,
	}
	if err := newToken.Tokenize(); err != nil {
		return nil, err
	}
	if err := newToken.Encrypt(); err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	newToken.Id = id
	newToken.CreatedAt = time.Now()
	newToken.UpdatedAt = time.Now()

	h.Store[newToken.Token] = newToken

	output := &NewTokenResponse{}
	output.Body.Token = newToken.Token

	return output, nil
}

type GetTokenRequest struct {
	Token string `path:"token" validate:"required"`
}

type GetTokenResponse struct {
	Body struct {
		Token models.Token `json:"encrypted_token"`
	}
}

func (h *BaseHandler) GetEncryptedToken(_ context.Context, in *GetTokenRequest) (*GetTokenResponse, error) {
	token := in.Token
	if token == "" {
		return nil, huma.Error400BadRequest("token is required")
	}

	stored := h.Store[token]
	stored.Payload = ""
	output := &GetTokenResponse{}
	output.Body.Token = stored
	return output, nil
}

func (h *BaseHandler) GetDecryptedToken(_ context.Context, in *GetTokenRequest) (*GetTokenResponse, error) {
	token := in.Token
	if token == "" {
		return nil, huma.Error400BadRequest("token is required")
	}

	stored := h.Store[token]
	payload, err := stored.Decrypt()
	if err != nil {
		return nil, err
	}
	stored.Payload = payload
	output := &GetTokenResponse{}
	output.Body.Token = stored
	return output, nil
}
