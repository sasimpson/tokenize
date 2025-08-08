package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
)

// Routes will register routes that are attached to the handler
func Routes() *mux.Router {
	r := mux.NewRouter()
	humaApi := humamux.New(r, huma.DefaultConfig("Tokenize", "3.0.0"))

	handlers := BaseHandler{}
	huma.AutoRegister(humaApi, handlers)

	return r
}
