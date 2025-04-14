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
