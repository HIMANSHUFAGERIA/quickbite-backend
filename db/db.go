package db

import (
	"context"
	"fmt"
	"log"

	"quickbite/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global connection pool
var DB *pgxpool.Pool

func Connect(cfg *config.Config) {

	var connStr string

	if cfg.DatabaseURL != "" {
		// ✅ Production (Railway)
		connStr = cfg.DatabaseURL
	} else {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
		)
	}
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Ping to verify connection is alive
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ Database unreachable: %v", err)
	}

	DB = pool
	log.Println("✅ Connected to PostgreSQL successfully")
}
