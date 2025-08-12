package api

import (
	"tokenize/persistence"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

type BaseHandler struct {
	Store persistence.Store
}

// Routes will register routes that are attached to the handler
func Routes(handlers *BaseHandler) *mux.Router {
	r := mux.NewRouter()
	humaApi := humamux.New(r, huma.DefaultConfig("Tokenize", "3.0.0"))

	huma.AutoRegister(humaApi, handlers)

	return r
}
