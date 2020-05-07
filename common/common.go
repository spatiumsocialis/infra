package common

import (
	"github.com/safe-distance/circle/pkg/models"
	"github.com/safe-distance/initialize"
)

// InitializeService loads the environment variables and connects to the database and automigrates
func InitializeService(overload bool, envFilenames ...string) {
	initialize.Env(overload, envFilenames...)
	db = initialize.DB(&User{}, &models.Circle{})
}
