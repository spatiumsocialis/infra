package common

import (
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from the provided files.
// If no filenames are passed, it loads variables from "./.env" by default.
func LoadEnv() error {
	rootDir := os.Getenv("PROJECT_ROOT")
	if err := godotenv.Load(rootDir + "/.env"); err != nil {
		return err
	}
	return nil
}

type environmentVar string

var (
	dbProvider         = environmentVar("DB_PROVIDER")
	dbConnectionString = environmentVar("DB_CONNECTION_STRING")
)

func (e environmentVar) String() string {
	return string(e)
}
