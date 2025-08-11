package api

import (
	"context"
	"net/http"

	"tokenize/models"

	"github.com/danielgtaylor/huma/v2"
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

	huma.Register(api, huma.Operation{
		OperationID:   "DeleteToken",
		Summary:       "Delete a token",
		Method:        http.MethodDelete,
		Path:          "/token/{token}",
		DefaultStatus: http.StatusNoContent,
		Errors: []int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusBadRequest,
			http.StatusNotFound,
		},
	}, h.DeleteToken)

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

func (h *BaseHandler) CreateToken(ctx context.Context, in *NewTokenRequest) (*NewTokenResponse, error) {

	newToken := models.Token{
		CreateToken: in.Body.Data,
	}
	if err := newToken.Tokenize(); err != nil {
		return nil, err
	}
	if err := newToken.Encrypt(); err != nil {
		return nil, err
	}

	tokenVal, err := h.Store.CreateToken(ctx, &newToken)
	if err != nil {
		return nil, err
	}

	output := &NewTokenResponse{}
	output.Body.Token = tokenVal.Token

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

func (h *BaseHandler) GetEncryptedToken(ctx context.Context, in *GetTokenRequest) (*GetTokenResponse, error) {
	token := in.Token
	if token == "" {
		return nil, huma.Error400BadRequest("token is required")
	}

	tokenVal, err := h.Store.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}

	tokenVal.Payload = ""
	output := &GetTokenResponse{}
	output.Body.Token = *tokenVal
	return output, nil
}

func (h *BaseHandler) GetDecryptedToken(ctx context.Context, in *GetTokenRequest) (*GetTokenResponse, error) {
	token := in.Token
	if token == "" {
		return nil, huma.Error400BadRequest("token is required")
	}

	tokenVal, err := h.Store.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}

	payload, err := tokenVal.Decrypt()
	if err != nil {
		return nil, err
	}

	tokenVal.Payload = payload
	output := &GetTokenResponse{}
	output.Body.Token = *tokenVal
	return output, nil
}

func (h *BaseHandler) DeleteToken(ctx context.Context, in *GetTokenRequest) (*struct{}, error) {
	token := in.Token

	tokenVal, err := h.Store.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}

	err = h.Store.DeleteToken(ctx, tokenVal)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
