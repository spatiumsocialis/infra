package common

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

// NewDB connects to the database and migrates it with the argument models.
func NewDB(models ...interface{}) *gorm.DB {
	// Get the DB_PROVIDER env var
	provider := os.Getenv(dbProvider.String())
	if provider == "" {
		log.Fatal("Error initializing DB: DB_PROVIDER env variable not set")
	}
	// Get the DB_CONNECTION_STRING env var
	connectionString := os.Getenv(dbConnectionString.String())
	if connectionString == "" {
		log.Fatal("Error initializing DB: DB_CONNECTION_STRING env variable not set")
	}

	// Open connection to the database
	log.Println("Connecting to db")
	db, err := gorm.Open(provider, connectionString)
	if err != nil {
		log.Fatalf("Error connecting to db: %s", err.Error())
	}

	// Migrate the schema
	log.Println("Auto-migrating schema")
	for _, model := range models {
		db = db.AutoMigrate(model)
	}
	return db
}
