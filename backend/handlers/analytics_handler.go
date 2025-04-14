package handlers

import (
	"context"
	"log"

	"digital-library/backend/database"

	"github.com/gofiber/fiber/v2"
)

// BorrowCount represents the structure for book borrow counts
type BorrowCount struct {
	BookID    int    `json:"book_id"`
	BookTitle string `json:"book_title"`
	Borrows   int    `json:"borrows"`
}

// GetMostBorrowedBooks retrieves books ordered by borrow count
func GetMostBorrowedBooks(c *fiber.Ctx) error {
	// Default limit to top 10, could make this a query param later
	limit := 10

	query := `SELECT 
	            b.id AS book_id, 
	            b.title AS book_title,
	            COUNT(lr.id) AS borrows
	          FROM books b
	          JOIN lending_records lr ON b.id = lr.book_id
	          GROUP BY b.id, b.title
	          ORDER BY borrows DESC
	          LIMIT $1`

	rows, err := database.DB.Query(context.Background(), query, limit)
	if err != nil {
		log.Printf("Error fetching most borrowed books: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve most borrowed books data",
		})
	}
	defer rows.Close()

	results := make([]BorrowCount, 0)
	for rows.Next() {
		var bc BorrowCount
		err := rows.Scan(&bc.BookID, &bc.BookTitle, &bc.Borrows)
		if err != nil {
			log.Printf("Error scanning borrow count row: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing borrow count data",
			})
		}
		results = append(results, bc)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating borrow count rows: %v", rows.Err())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving borrow count data",
		})
	}

	return c.JSON(results)
}

// MonthlyTrend represents the structure for monthly lending counts
type MonthlyTrend struct {
	Month string `json:"month"` // Format YYYY-MM
	Count int    `json:"count"`
}

// GetMonthlyLendingTrends retrieves lending counts grouped by month
func GetMonthlyLendingTrends(c *fiber.Ctx) error {
	// Query to count records per month based on borrow_date
	// Adjust the time period (e.g., last 12 months) if needed
	query := `SELECT 
	            to_char(borrow_date, 'YYYY-MM') AS month, 
	            COUNT(*) AS count
	          FROM lending_records
	          -- WHERE borrow_date >= NOW() - INTERVAL '12 months' -- Optional: Filter by time period
	          GROUP BY month
	          ORDER BY month ASC`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error fetching monthly lending trends: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve monthly lending trends data",
		})
	}
	defer rows.Close()

	results := make([]MonthlyTrend, 0)
	for rows.Next() {
		var mt MonthlyTrend
		err := rows.Scan(&mt.Month, &mt.Count)
		if err != nil {
			log.Printf("Error scanning monthly trend row: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing monthly trend data",
			})
		}
		results = append(results, mt)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating monthly trend rows: %v", rows.Err())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving monthly trend data",
		})
	}

	return c.JSON(results)
}

// CategoryDistribution represents the count of books per category
type CategoryDistribution struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// GetCategoryDistribution retrieves the count of books per category
func GetCategoryDistribution(c *fiber.Ctx) error {
	// Query to count books per category. Handle NULL categories.
	query := `SELECT 
	            COALESCE(category, 'Uncategorized') AS category,
	            COUNT(*) AS count
	          FROM books
	          GROUP BY COALESCE(category, 'Uncategorized')
	          ORDER BY count DESC`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error fetching category distribution: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve category distribution data",
		})
	}
	defer rows.Close()

	results := make([]CategoryDistribution, 0)
	for rows.Next() {
		var cd CategoryDistribution
		err := rows.Scan(&cd.Category, &cd.Count)
		if err != nil {
			log.Printf("Error scanning category distribution row: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing category distribution data",
			})
		}
		results = append(results, cd)
	}

	if rows.Err() != nil {
		log.Printf("Error iterating category distribution rows: %v", rows.Err())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving category distribution data",
		})
	}

	return c.JSON(results)
}
