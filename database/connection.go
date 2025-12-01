package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"backenduas/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *pgxpool.Pool
var MongoClient *mongo.Client
var MongoDB *mongo.Database

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

// ===============================
// CONNECT MONGODB
// ===============================
func ConnectMongo() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    uri := config.AppEnv.MongoURI
    dbName := config.AppEnv.MongoDB

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("❌ Failed to connect MongoDB: %v", err)
    }

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("❌ Failed to ping MongoDB: %v", err)
    }

    MongoClient = client
    MongoDB = client.Database(dbName)

    fmt.Println("🍃 MongoDB connected:", dbName)
}

// ===============================
// INIT BOTH DATABASES
// ===============================
func ConnectDatabases() {
	ConnectPostgre()
	ConnectMongo()
}
