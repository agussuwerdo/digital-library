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

	// Optional: Configure pool settings
	dbConfig.MaxConns = 10                      // Example: Set max connections
	dbConfig.MinConns = 2                       // Example: Set min connections
	dbConfig.MaxConnLifetime = time.Hour        // Example: Max connection lifetime
	dbConfig.MaxConnIdleTime = time.Minute * 30 // Example: Max idle time

	// Connect to the database
	DB, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Optional: Ping the database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
