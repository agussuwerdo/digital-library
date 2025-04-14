// handler/handler.go
package handler

import (
	"net/http"

	"digital-library/backend/app"
	"digital-library/backend/database"

	"github.com/gofiber/adaptor/v2"
)

// Handler for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	app := app.SetupApp()
	defer database.Close()

	handler := adaptor.FiberApp(app)
	handler.ServeHTTP(w, r)
}
