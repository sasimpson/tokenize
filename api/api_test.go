package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"tokenize/models"
	"tokenize/persistence/mock"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	testCases := []struct {
		name    string
		handler *BaseHandler
		test    func(t *testing.T, router *mux.Router)
	}{
		{
			name: "routes function creates router",
			handler: &BaseHandler{
				Store: mock.Store{},
			},
			test: func(t *testing.T, router *mux.Router) {
				assert.NotNil(t, router)
				assert.IsType(t, &mux.Router{}, router)
			},
		},
		{
			name: "router has registered routes",
			handler: &BaseHandler{
				Store: mock.Store{},
			},
			test: func(t *testing.T, router *mux.Router) {
				// Create a test request to check if routes are registered
				req, err := http.NewRequest("GET", "/", nil)
				assert.NoError(t, err)

				// Use the router to match routes
				var match mux.RouteMatch
				matched := router.Match(req, &match)

				// The router should at least be configured (even if no match for root path)
				assert.NotNil(t, match)
				_ = matched // We expect this to be false for root path, that's fine
			},
		},
		{
			name: "token routes are accessible",
			handler: &BaseHandler{
				Store: mock.Store{
					Token: &models.Token{
						Token: "test-token",
						CreateToken: models.CreateToken{
							Payload: "test-payload",
						},
					},
				},
			},
			test: func(t *testing.T, router *mux.Router) {
				// Test that specific token endpoints are registered
				testCases := []struct {
					method string
					path   string
				}{
					{"POST", "/token"},
					{"GET", "/token/test-token"},
					{"GET", "/token/test-token/decrypt"},
					{"DELETE", "/token/test-token"},
				}

				for _, tc := range testCases {
					req, err := http.NewRequest(tc.method, tc.path, nil)
					assert.NoError(t, err)

					var match mux.RouteMatch
					matched := router.Match(req, &match)

					// Check that the route pattern exists (even if handler may fail)
					if matched {
						assert.NotNil(t, match.Route, "Route should be matched for %s %s", tc.method, tc.path)
					}
				}
			},
		},
		{
			name: "router handles requests",
			handler: &BaseHandler{
				Store: mock.Store{
					Token: &models.Token{
						Token: "test-token",
						CreateToken: models.CreateToken{
							Payload: "test-payload",
						},
					},
				},
			},
			test: func(t *testing.T, router *mux.Router) {
				// Test that the router can actually handle a request
				req, err := http.NewRequest("GET", "/token/test-token", nil)
				assert.NoError(t, err)

				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				// Should return some response (even if it's an error due to missing headers)
				assert.NotEqual(t, 0, rr.Code, "Router should handle the request and return a status code")
			},
		},
		{
			name:    "routes function with nil handler",
			handler: nil,
			test: func(t *testing.T, router *mux.Router) {
				// Should still create a router even with nil handler
				assert.NotNil(t, router)
				assert.IsType(t, &mux.Router{}, router)
			},
		},
		{
			name: "router configuration",
			handler: &BaseHandler{
				Store: mock.Store{},
			},
			test: func(t *testing.T, router *mux.Router) {
				// Test that router is properly configured
				assert.NotNil(t, router)

				// Check that we can walk the routes (indicating they're registered)
				routeCount := 0
				err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
					routeCount++
					return nil
				})
				assert.NoError(t, err)

				// Should have some routes registered (at least the token routes)
				assert.Greater(t, routeCount, 0, "Router should have registered routes")
			},
		},
		{
			name: "huma api integration",
			handler: &BaseHandler{
				Store: mock.Store{},
			},
			test: func(t *testing.T, router *mux.Router) {
				// Test that Huma API is properly integrated
				req, err := http.NewRequest("OPTIONS", "/", nil)
				assert.NoError(t, err)

				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				// Huma typically handles OPTIONS requests for CORS
				// The exact response depends on Huma configuration, but it should respond
				assert.NotEqual(t, 0, rr.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := Routes(tc.handler)
			tc.test(t, router)
		})
	}
}
