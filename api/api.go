package api

import "tokenize/models"

type BaseHandler struct {
	Store map[string]models.Token
}
