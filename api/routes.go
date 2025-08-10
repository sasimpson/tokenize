package api

import (
	"tokenize/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

// Routes will register routes that are attached to the handler
func Routes() *mux.Router {
	r := mux.NewRouter()
	humaApi := humamux.New(r, huma.DefaultConfig("Tokenize", "3.0.0"))

	store := make(map[string]models.Token)
	handlers := &BaseHandler{
		Store: store,
	}
	huma.AutoRegister(humaApi, handlers)

	return r
}
