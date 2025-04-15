// handler/handler.go
package handler

import (
	"net/http"
	"strings"

	"digital-library/backend/app"
	"digital-library/backend/database"

	"github.com/gofiber/adaptor/v2"
)

// Handler for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	// Create a new Fiber app for each request
	app := app.SetupApp()
	defer database.Close()

	// Create a custom handler that preserves the original request
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new request with the original path and query
		newReq := r.Clone(r.Context())

		// Ensure the path starts with /api
		if !strings.HasPrefix(r.URL.Path, "/api") {
			newReq.URL.Path = "/api" + r.URL.Path
		}

		// Create the full URI with query parameters
		fullURI := newReq.URL.Path
		if r.URL.RawQuery != "" {
			fullURI += "?" + r.URL.RawQuery
		}
		newReq.RequestURI = fullURI

		// Convert to Fiber handler and serve
		fiberHandler := adaptor.FiberApp(app)
		fiberHandler.ServeHTTP(w, newReq)
	})

	handler.ServeHTTP(w, r)
}
