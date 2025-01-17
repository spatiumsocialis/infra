package common

import (
	"errors"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // gorm driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // gorm driver
)

// NewDB connects to the database and migrates it with the argument models.
func NewDB(models ...interface{}) (*gorm.DB, error) {
	// Get the DB_PROVIDER env var
	provider := os.Getenv(dbProvider.String())
	if provider == "" {
		return nil, errors.New("Error connecting to DB: DB_PROVIDER env variable not set")
	}
	// Get the DB_CONNECTION_STRING env var
	connectionString := os.Getenv(dbConnectionString.String())
	if connectionString == "" {
		return nil, errors.New("Error connecting to DB: DB_CONNECTION_STRING env variable not set")
	}

	// Open connection to the database
	log.Println("Connecting to db")
	db, err := gorm.Open(provider, connectionString)
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	log.Println("Auto-migrating schema")
	if err := db.AutoMigrate(models...).Error; err != nil {
		log.Fatalf("error migrating db: %v\n", err.Error())
	}

	return db, nil
}
