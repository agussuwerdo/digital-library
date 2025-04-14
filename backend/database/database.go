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

	// Optimize for serverless environment
	dbConfig.MaxConns = 5                        // Reduced max connections for serverless
	dbConfig.MinConns = 1                        // Keep minimum connections low
	dbConfig.MaxConnLifetime = time.Minute * 5   // Shorter lifetime for serverless
	dbConfig.MaxConnIdleTime = time.Minute * 1   // Shorter idle time
	dbConfig.HealthCheckPeriod = time.Second * 5 // More frequent health checks

	// Add connection retry logic
	var pool *pgxpool.Pool
	for i := 0; i < 3; i++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
		if err == nil {
			break
		}
		log.Printf("Attempt %d: Failed to connect to database: %v\n", i+1, err)
		time.Sleep(time.Second * 2)
	}

	if err != nil {
		log.Fatalf("Unable to connect to database after retries: %v\n", err)
	}

	DB = pool

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
