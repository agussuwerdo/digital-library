package models

import (
	"time"
)

// Book represents the structure for a book in the library
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	Quantity  int       `json:"quantity"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LendingRecord represents the structure for a lending record
type LendingRecord struct {
	ID         int        `json:"id"`
	BookID     int        `json:"book_id"` // Foreign key to Book
	Borrower   string     `json:"borrower"`
	BorrowDate time.Time  `json:"borrow_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"` // Pointer to allow null
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// MonthlyTrend represents the structure for monthly lending counts
type MonthlyTrend struct {
	Month string `json:"month"` // Format YYYY-MM
	Count int    `json:"count"`
}

// CategoryDistribution represents the count of books per category
type CategoryDistribution struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
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

// LoginRequest represents the structure for login requests
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the structure for login responses
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// RegisterRequest represents the structure for registration requests
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// RegisterResponse represents the structure for registration responses
type RegisterResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
