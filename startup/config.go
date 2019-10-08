package startup

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Config *Configuration

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	Config = &Configuration{
		Database: Database{
			UserName: os.Getenv("MYSQL_USER"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			DbName:   os.Getenv("MYSQL_DATABASE"),
		},
		GoogleApiKey: os.Getenv("GOOGLE_API_KEY"),
	}
}

type Configuration struct {
	Database     Database
	GoogleApiKey string
}

type Database struct {
	UserName string
	Password string
	DbName   string
}
