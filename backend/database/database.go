package database

import (
	"context"
	"log"
	"os"
	"time"

	"digital-library/backend/config" // Import the config package

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Connect initializes the database connection pool
func Connect(cfg *config.Config) {
	var err error
	// Use pgxpool.ParseConfig to handle the DSN
	dbConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v\n", err)
	}

	// Configure pool settings for production
	dbConfig.MaxConns = 25                      // Increased max connections
	dbConfig.MinConns = 5                       // Increased min connections
	dbConfig.MaxConnLifetime = time.Hour * 2    // Increased connection lifetime
	dbConfig.MaxConnIdleTime = time.Minute * 30 // Keep idle time the same
	dbConfig.ConnConfig.RuntimeParams = map[string]string{
		"application_name": "digital_library",
		"search_path":      "public",
	}

	// Connect to the database with retry logic
	var pool *pgxpool.Pool
	for i := 0; i < 3; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
		if err == nil {
			break
		}
		log.Printf("Attempt %d to connect to database failed: %v\n", i+1, err)
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Fatalf("Unable to connect to database after retries: %v\n", err)
	}
	DB = pool

	// Verify connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = DB.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Database connection established successfully.")
}

// Close closes the database connection pool
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed.")
	}
}

// --- Helper function for graceful shutdown ---
func SetupCloseHandler() {
	c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM) // syscall requires import
	// TODO: Add proper signal handling if needed
	go func() {
		<-c
		log.Println("\r- Ctrl+C pressed in Terminal")
		Close()
		os.Exit(0)
	}()
}
