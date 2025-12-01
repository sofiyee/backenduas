package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	AppPort string
	DBDsn   string
	ApiKey  string
	MongoURI string `mapstructure:"MONGO_URI"`
    MongoDB  string `mapstructure:"MONGO_DB"`
}

var AppEnv *Env

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	AppEnv = &Env{
		AppPort: os.Getenv("APP_PORT"),
		DBDsn:   os.Getenv("DB_DSN"),
		ApiKey:  os.Getenv("API_KEY"),
		MongoURI: os.Getenv("MONGO_URI"),
		MongoDB:  os.Getenv("MONGO_DB"),


	}
}
