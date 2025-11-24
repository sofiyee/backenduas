package database

import (
	"context"
	"log"
	"time"

	"backenduas/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectPostgre() {
	dsn := config.AppEnv.DBDsn

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("❌ Failed to create DB pool: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	DB = db
	log.Println("✅ PostgreSQL connected successfully!")
}
