package api

import (
	"net/http"

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
	}, nil)

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
	}, nil)

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
	}, nil)

}
