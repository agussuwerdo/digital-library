package handlers

import (
	"time"

	"digital-library/backend/config" // Assuming config is setup to load JWTSecret

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// LoginPayload defines the expected structure for the login request body
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login handles user authentication and JWT generation
func Login(cfg *config.Config) fiber.Handler { // Pass config to access JWT secret
	return func(c *fiber.Ctx) error {
		payload := new(LoginPayload)
		if err := c.BodyParser(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// --- Hardcoded user validation (Replace with DB lookup in a real app) ---
		// IMPORTANT: This is insecure and only for demonstration purposes.
		if payload.Username != "librarian" || payload.Password != "password123" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// --- Generate JWT ---
		claims := jwt.MapClaims{
			"username": payload.Username,
			"role":     "librarian",                           // Example claim
			"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
			"iat":      time.Now().Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate token",
			})
		}

		return c.JSON(fiber.Map{"token": t})
	}
}
