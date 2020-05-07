package common

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads the environment variables from the provided files.
// If no filenames are passed, it loads variables from "./.env" by default.
func LoadEnv(overload bool, envFilenames ...string) {
	var err error
	if overload {
		err = godotenv.Overload(envFilenames...)
	}
	err = godotenv.Load(envFilenames...)
	if err != nil && len(envFilenames) > 0 {
		log.Fatalf("Error loading environment files: %v\n", err.Error())
	}
}

type environmentVar string

var (
	dbProvider         = environmentVar("DB_PROVIDER")
	dbConnectionString = environmentVar("DB_CONNECTION_STRING")
)

func (e environmentVar) String() string {
	return string(e)
}
