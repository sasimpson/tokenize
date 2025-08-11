package api

import (
	"tokenize/persistence"
)

type BaseHandler struct {
	Store persistence.Store
}
