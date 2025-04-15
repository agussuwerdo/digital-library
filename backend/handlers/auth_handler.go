package handlers

import (
	"context"
	"log"
	"strings"
	"time"

	"digital-library/backend/config"
	"digital-library/backend/database"
	"digital-library/backend/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// @Summary Register a new user
// @Description Create a new user account with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func Register(c *fiber.Ctx) error {
	payload := new(models.RegisterRequest)
	if err := c.BodyParser(payload); err != nil {
		log.Printf("Error parsing registration payload: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Basic validation
	if payload.Username == "" || payload.Password == "" || payload.Email == "" {
		log.Printf("Invalid registration attempt - missing fields. Username: %v, Email: %v",
			payload.Username != "", payload.Email != "")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, password, and email are required",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not hash password",
		})
	}

	// Insert the new user
	query := `INSERT INTO users (username, password_hash, email, role) 
	          VALUES ($1, $2, $3, 'user') 
	          RETURNING id, username, email, role, created_at, updated_at`

	var user models.User
	err = database.DB.QueryRow(context.Background(), query,
		payload.Username, string(hashedPassword), payload.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_username_key\"") {
			log.Printf("Registration failed: Duplicate username '%s'", payload.Username)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		} else if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			log.Printf("Registration failed: Duplicate email '%s'", payload.Email)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		log.Printf("Database error during registration: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Could not create user. Please try again later.",
			"details": err.Error(),
		})
	}

	log.Printf("User registered successfully: %s (ID: %d)", user.Username, user.ID)
	return c.Status(fiber.StatusCreated).JSON(user)
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func Login(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := new(models.LoginRequest)
		if err := c.BodyParser(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// Query the user from the database
		var user models.User
		var passwordHash string
		query := `SELECT id, username, email, role, password_hash, created_at, updated_at 
		          FROM users WHERE username = $1 OR email = $1`

		err := database.DB.QueryRow(context.Background(), query, payload.Username).
			Scan(&user.ID, &user.Username, &user.Email, &user.Role, &passwordHash, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// Compare the provided password with the stored hash
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(payload.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// Generate JWT
		claims := jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"role":     user.Role,
			"exp":      time.Now().Add(time.Hour * 72).Unix(),
			"iat":      time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate token",
			})
		}

		return c.JSON(fiber.Map{
			"token": t,
			"user":  user,
		})
	}
}
