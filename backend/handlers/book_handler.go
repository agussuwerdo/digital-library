package handlers

import (
	"context"
	"log"

	// Needed for models.Book
	"digital-library/backend/database"
	"digital-library/backend/models"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v5" // Not needed here
	"github.com/jackc/pgx/v5"
)

// --- Book Handlers ---

// CreateBook handles the creation of a new book
func CreateBook(c *fiber.Ctx) error {
	book := new(models.Book)

	// Parse the request body into the book struct
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Basic validation (Add more as needed)
	if book.Title == "" || book.Author == "" || book.ISBN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title, Author, and ISBN are required fields",
		})
	}
	if book.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Quantity cannot be negative",
		})
	}

	// SQL query to insert the book
	query := `INSERT INTO books (title, author, isbn, quantity, category) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, created_at, updated_at`

	// Execute the query
	row := database.DB.QueryRow(context.Background(), query,
		book.Title, book.Author, book.ISBN, book.Quantity, book.Category)

	// Scan the returned id, created_at, updated_at into the book struct
	err := row.Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		// Handle potential errors, e.g., duplicate ISBN
		log.Printf("Error creating book: %v", err) // Log the error
		// TODO: Check for specific pgx errors like unique constraint violation (e.g., using pgconn.PgError)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create book", // Consider more specific error messages
		})
	}

	// Return the newly created book
	return c.Status(fiber.StatusCreated).JSON(book)
}

// GetBooks retrieves all books from the database
func GetBooks(c *fiber.Ctx) error {
	query := `SELECT id, title, author, isbn, quantity, category, created_at, updated_at 
	          FROM books ORDER BY title ASC`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error fetching books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve books",
		})
	}
	defer rows.Close()

	books := make([]models.Book, 0)

	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN,
			&book.Quantity, &book.Category, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning book row: %v", err)
			// Decide if you want to return partial results or an error
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing book data",
			})
		}
		books = append(books, book)
	}

	if rows.Err() != nil { // Check for errors during iteration
		log.Printf("Error iterating book rows: %v", rows.Err())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving book data",
		})
	}

	return c.JSON(books)
}

// GetBook retrieves a single book by its ID
func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Add validation to ensure id is a number if needed

	query := `SELECT id, title, author, isbn, quantity, category, created_at, updated_at 
	          FROM books WHERE id = $1`

	var book models.Book
	row := database.DB.QueryRow(context.Background(), query, id)

	err := row.Scan(
		&book.ID, &book.Title, &book.Author, &book.ISBN,
		&book.Quantity, &book.Category, &book.CreatedAt, &book.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Book not found
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		// Other potential errors
		log.Printf("Error fetching book %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve book",
		})
	}

	return c.JSON(book)
}

// UpdateBook handles updating an existing book
func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Add validation for ID

	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Basic validation
	if book.Title == "" || book.Author == "" || book.ISBN == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title, Author, and ISBN are required fields",
		})
	}
	if book.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Quantity cannot be negative",
		})
	}

	query := `UPDATE books 
	          SET title = $1, author = $2, isbn = $3, quantity = $4, category = $5, updated_at = NOW() 
	          WHERE id = $6
	          RETURNING id, title, author, isbn, quantity, category, created_at, updated_at`

	var updatedBook models.Book
	row := database.DB.QueryRow(context.Background(), query,
		book.Title, book.Author, book.ISBN, book.Quantity, book.Category, id)

	err := row.Scan(
		&updatedBook.ID, &updatedBook.Title, &updatedBook.Author, &updatedBook.ISBN,
		&updatedBook.Quantity, &updatedBook.Category, &updatedBook.CreatedAt, &updatedBook.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Book not found to update
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		// Handle other errors like unique constraint violation on ISBN
		log.Printf("Error updating book %s: %v", id, err)
		// TODO: Check for specific pgx errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update book",
		})
	}

	return c.JSON(updatedBook)
}

// DeleteBook handles deleting a book by its ID
func DeleteBook(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Add validation for ID

	query := `DELETE FROM books WHERE id = $1 RETURNING id` // RETURNING id to check if deletion happened

	var deletedID int
	row := database.DB.QueryRow(context.Background(), query, id)
	err := row.Scan(&deletedID)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Book not found to delete
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		// Other potential errors (e.g., foreign key constraints if ON DELETE CASCADE wasn't set correctly, though we set it)
		log.Printf("Error deleting book %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete book",
		})
	}

	// Return success message or no content
	// return c.SendStatus(fiber.StatusNoContent) // Option 1: No content
	return c.JSON(fiber.Map{"message": "Book deleted successfully", "id": deletedID}) // Option 2: Confirmation message
}
