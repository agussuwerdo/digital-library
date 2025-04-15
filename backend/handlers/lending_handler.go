package handlers

import (
	"context"
	"log"
	"strconv"
	"time"

	// For custom errors
	"digital-library/backend/database"
	"digital-library/backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	// For checking specific error codes
)

// LendBookPayload defines the expected structure for the lend book request
type LendBookPayload struct {
	BookID   int    `json:"book_id"`
	Borrower string `json:"borrower"`
}

// @Summary Lend a book
// @Description Create a new lending record for a book
// @Tags lending
// @Accept json
// @Produce json
// @Param lending body models.LendingRecord true "Lending record object"
// @Success 201 {object} models.LendingRecord
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lending [post]
func LendBook(c *fiber.Ctx) error {
	payload := new(LendBookPayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Basic validation
	if payload.BookID <= 0 || payload.Borrower == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Book ID and Borrower name are required"})
	}

	// Use a transaction to ensure atomicity
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not start transaction"})
	}
	// Defer rollback in case of error, commit will override this if successful
	defer tx.Rollback(context.Background())

	// 1. Check book quantity and lock the row for update
	var currentQuantity int
	checkQuery := `SELECT quantity FROM books WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(context.Background(), checkQuery, payload.BookID).Scan(&currentQuantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
		}
		log.Printf("Error checking book quantity: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not verify book availability"})
	}

	if currentQuantity <= 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Book is currently out of stock"})
	}

	// 2. Decrease book quantity
	updateQuery := `UPDATE books SET quantity = quantity - 1, updated_at = NOW() WHERE id = $1`
	_, err = tx.Exec(context.Background(), updateQuery, payload.BookID)
	if err != nil {
		log.Printf("Error updating book quantity: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update book quantity"})
	}

	// 3. Create lending record
	insertQuery := `INSERT INTO lending_records (book_id, borrower_name, borrow_date) 
	                 VALUES ($1, $2, $3) 
	                 RETURNING id, borrow_date, created_at, updated_at`
	var newRecord models.LendingRecord
	newRecord.BookID = payload.BookID // Populate from payload
	newRecord.Borrower = payload.Borrower
	borrowDate := time.Now().UTC().Truncate(24 * time.Hour) // Use UTC date part only

	row := tx.QueryRow(context.Background(), insertQuery,
		payload.BookID, payload.Borrower, borrowDate)
	err = row.Scan(&newRecord.ID, &newRecord.BorrowDate, &newRecord.CreatedAt, &newRecord.UpdatedAt)
	if err != nil {
		log.Printf("Error creating lending record: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create lending record"})
	}

	// 4. Commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not complete lending operation"})
	}

	return c.Status(fiber.StatusCreated).JSON(newRecord)
}

// @Summary Return a book
// @Description Mark a lending record as returned and update book availability
// @Tags lending
// @Accept json
// @Produce json
// @Param id path int true "Lending Record ID"
// @Success 200 {object} models.LendingRecord
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lending/{id}/return [put]
func ReturnBook(c *fiber.Ctx) error {
	lendingRecordID := c.Params("id")
	// TODO: Add validation for lendingRecordID

	returnDate := time.Now().UTC().Truncate(24 * time.Hour) // Use UTC date part only

	// Use a transaction
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not start transaction"})
	}
	defer tx.Rollback(context.Background())

	// 1. Update lending record and get the book_id
	var bookID int
	updateLendingQuery := `UPDATE lending_records 
	                       SET return_date = $1, updated_at = NOW() 
	                       WHERE id = $2 AND return_date IS NULL -- Only update if not already returned
	                       RETURNING book_id`
	err = tx.QueryRow(context.Background(), updateLendingQuery, returnDate, lendingRecordID).Scan(&bookID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Either record doesn't exist or was already returned
			// Check if record exists but is already returned
			var exists bool
			checkExistsQuery := `SELECT EXISTS(SELECT 1 FROM lending_records WHERE id = $1 AND return_date IS NOT NULL)`
			errCheck := database.DB.QueryRow(context.Background(), checkExistsQuery, lendingRecordID).Scan(&exists)
			if errCheck == nil && exists {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Book already returned"})
			}
			// Otherwise, the record wasn't found
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lending record not found or already returned"})
		}
		// Other errors
		log.Printf("Error updating lending record %s: %v", lendingRecordID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update lending record"})
	}

	// 2. Increment book quantity
	updateBookQuery := `UPDATE books SET quantity = quantity + 1, updated_at = NOW() WHERE id = $1`
	_, err = tx.Exec(context.Background(), updateBookQuery, bookID)
	if err != nil {
		// This shouldn't ideally happen if the foreign key constraint is working, but handle it
		log.Printf("Error incrementing book quantity for book %d: %v", bookID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update book quantity after return"})
	}

	// 3. Commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not complete return operation"})
	}

	return c.JSON(fiber.Map{"message": "Book returned successfully"})
}

// LendingRecordDetail extends LendingRecord to include book details
type LendingRecordDetail struct {
	models.LendingRecord        // Embed LendingRecord
	BookTitle            string `json:"book_title"`
	BookAuthor           string `json:"book_author"`
}

// @Summary Get lending records
// @Description Get all lending records with optional search and filtering
// @Tags lending
// @Accept json
// @Produce json
// @Param search query string false "Search term"
// @Param status query string false "Filter by status (active/returned)"
// @Param book_id query int false "Filter by book ID"
// @Success 200 {array} models.LendingRecord
// @Failure 500 {object} map[string]string
// @Router /lending [get]
func GetLendingRecords(c *fiber.Ctx) error {
	// Get query parameters
	search := c.Query("search", "")
	borrower := c.Query("borrower", "")
	status := c.Query("status", "") // "active" or "returned"
	bookTitle := c.Query("bookTitle", "")

	// Build the base query
	query := `SELECT 
	            lr.id, lr.book_id, lr.borrower_name, lr.borrow_date, lr.return_date, 
	            lr.created_at, lr.updated_at, 
	            b.title AS book_title, b.author AS book_author
	          FROM lending_records lr
	          JOIN books b ON lr.book_id = b.id
	          WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	// Add search condition (matches borrower name or book title)
	if search != "" {
		query += ` AND (LOWER(lr.borrower_name) LIKE LOWER($` + strconv.Itoa(argCount) + `) OR LOWER(b.title) LIKE LOWER($` + strconv.Itoa(argCount) + `))`
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Add borrower filter
	if borrower != "" {
		query += ` AND LOWER(lr.borrower_name) = LOWER($` + strconv.Itoa(argCount) + `)`
		args = append(args, borrower)
		argCount++
	}

	// Add status filter
	if status == "active" {
		query += ` AND lr.return_date IS NULL`
	} else if status == "returned" {
		query += ` AND lr.return_date IS NOT NULL`
	}

	// Add book title filter
	if bookTitle != "" {
		query += ` AND LOWER(b.title) = LOWER($` + strconv.Itoa(argCount) + `)`
		args = append(args, bookTitle)
		argCount++
	}

	// Add sorting
	query += ` ORDER BY lr.borrow_date DESC, lr.created_at DESC`

	rows, err := database.DB.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Error fetching lending records: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve lending records",
		})
	}
	defer rows.Close()

	records := make([]LendingRecordDetail, 0)

	for rows.Next() {
		var record LendingRecordDetail
		err := rows.Scan(
			&record.ID, &record.BookID, &record.Borrower, &record.BorrowDate, &record.ReturnDate,
			&record.CreatedAt, &record.UpdatedAt,
			&record.BookTitle, &record.BookAuthor,
		)
		if err != nil {
			log.Printf("Error scanning lending record row: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing lending record data",
			})
		}
		records = append(records, record)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating lending record rows: %v", rows.Err())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving lending record data",
		})
	}

	return c.JSON(records)
}

// @Summary Delete a lending record
// @Description Delete a lending record by its ID
// @Tags lending
// @Accept json
// @Produce json
// @Param id path int true "Lending Record ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /lending/{id} [delete]
func DeleteLendingRecord(c *fiber.Ctx) error {
	lendingRecordID := c.Params("id")
	// TODO: Add validation for lendingRecordID

	// Use a transaction
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not start transaction"})
	}
	defer tx.Rollback(context.Background())

	// 1. Get book_id and check if the book was returned before deleting
	var bookID int
	var returnDate *time.Time
	selectQuery := `SELECT book_id, return_date FROM lending_records WHERE id = $1`
	err = tx.QueryRow(context.Background(), selectQuery, lendingRecordID).Scan(&bookID, &returnDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lending record not found"})
		}
		log.Printf("Error checking lending record %s before delete: %v", lendingRecordID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve lending record details"})
	}

	// 2. Delete the lending record
	deleteQuery := `DELETE FROM lending_records WHERE id = $1`
	commandTag, err := tx.Exec(context.Background(), deleteQuery, lendingRecordID)
	if err != nil {
		log.Printf("Error deleting lending record %s: %v", lendingRecordID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete lending record"})
	}
	if commandTag.RowsAffected() == 0 {
		// Should have been caught by the select earlier, but as a safeguard
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lending record not found"})
	}

	// 3. If the book was *not* returned, increment the book quantity back
	if returnDate == nil {
		updateBookQuery := `UPDATE books SET quantity = quantity + 1, updated_at = NOW() WHERE id = $1`
		_, err = tx.Exec(context.Background(), updateBookQuery, bookID)
		if err != nil {
			// Handle potential error updating book quantity
			log.Printf("Error incrementing book quantity for book %d after deleting unreturned record %s: %v", bookID, lendingRecordID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update book quantity after deleting lending record"})
		}
	}

	// 4. Commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not complete delete operation"})
	}

	return c.JSON(fiber.Map{"message": "Lending record deleted successfully"})
}
