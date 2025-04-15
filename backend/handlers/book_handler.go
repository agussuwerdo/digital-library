package handlers

import (
	"context"
	"log"
	"strconv"
	"strings"

	// Needed for models.Book
	"digital-library/backend/database"
	"digital-library/backend/models"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v5" // Not needed here
	"github.com/jackc/pgx/v5"
)

// @Summary Get all books
// @Description Get all books with optional search and filtering
// @Tags books
// @Accept json
// @Produce json
// @Param search query string false "Search term for title or author"
// @Param category query string false "Filter by category"
// @Param author query string false "Filter by author"
// @Param available query string false "Filter by availability (true/false)"
// @Success 200 {array} models.Book
// @Failure 500 {object} map[string]string
// @Router /books [get]
func GetBooks(c *fiber.Ctx) error {
	// Get query parameters
	search := c.Query("search", "")
	category := c.Query("category", "")
	author := c.Query("author", "")
	available := c.Query("available", "")

	// Build the base query
	query := `SELECT id, title, author, isbn, quantity, category, created_at, updated_at 
	          FROM books WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	// Add search condition (matches title or author)
	if search != "" {
		query += ` AND (LOWER(title) LIKE LOWER($` + strconv.Itoa(argCount) + `) OR LOWER(author) LIKE LOWER($` + strconv.Itoa(argCount) + `))`
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Add category filter
	if category != "" {
		query += ` AND LOWER(category) = LOWER($` + strconv.Itoa(argCount) + `)`
		args = append(args, category)
		argCount++
	}

	// Add author filter
	if author != "" {
		query += ` AND LOWER(author) = LOWER($` + strconv.Itoa(argCount) + `)`
		args = append(args, author)
		argCount++
	}

	// Add availability filter
	if available == "true" {
		query += ` AND quantity > 0`
	} else if available == "false" {
		query += ` AND quantity = 0`
	}

	// Add sorting
	query += ` ORDER BY title ASC`

	rows, err := database.DB.Query(context.Background(), query, args...)
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing book data",
			})
		}
		books = append(books, book)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating book rows: %v", rows.Err())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving book data",
		})
	}

	return c.JSON(books)
}

// @Summary Get a book by ID
// @Description Get a single book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [get]
func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")

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
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Book not found",
			})
		}
		log.Printf("Error fetching book %s: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve book",
		})
	}

	return c.JSON(book)
}

// @Summary Create a new book
// @Description Create a new book with the provided information
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book object"
// @Success 201 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func CreateBook(c *fiber.Ctx) error {
	book := new(models.Book)

	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

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

	query := `INSERT INTO books (title, author, isbn, quantity, category) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, created_at, updated_at`

	row := database.DB.QueryRow(context.Background(), query,
		book.Title, book.Author, book.ISBN, book.Quantity, book.Category)

	err := row.Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		log.Printf("Error creating book: %v", err)
		// Check if error is due to duplicate ISBN
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"books_isbn_key\"") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "A book with this ISBN already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create book",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

// @Summary Update a book
// @Description Update an existing book with the provided information
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body models.Book true "Book object"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [put]
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
		// Check if error is due to duplicate ISBN
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"books_isbn_key\"") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "A book with this ISBN already exists",
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

// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [delete]
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
