package api

import (
	"tokenize/api/tokens"
)

type BaseHandler struct {
	TokenHandler tokens.Handler
}
