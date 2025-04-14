package handlers

import (
	"context"
	"log"
	"time"

	"digital-library/backend/config"
	"digital-library/backend/database"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// LoginPayload defines the expected structure for the login request body
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterPayload defines the expected structure for the registration request body
type RegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Register handles user registration
func Register(c *fiber.Ctx) error {
	payload := new(RegisterPayload)
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

	var user User
	err = database.DB.QueryRow(context.Background(), query,
		payload.Username, string(hashedPassword), payload.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
			log.Printf("Registration failed: Duplicate username '%s'", payload.Username)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Username already exists",
			})
		} else if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
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

// Login handles user authentication and JWT generation
func Login(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := new(LoginPayload)
		if err := c.BodyParser(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// Query the user from the database
		var user User
		var passwordHash string
		query := `SELECT id, username, email, role, password_hash, created_at, updated_at 
		          FROM users WHERE username = $1`

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
