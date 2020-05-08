package common

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
)

// Service represents an HTTP service
type Service struct {
	Name       string
	PathPrefix string
	DB         *gorm.DB
}

// NewService returns a new service constructed with the given arguments
func NewService(name string, pathPrefix string, models ...interface{}) *Service {
	db, err := NewDB(models...)
	if err != nil {
		log.Fatalf("Error initializing DB: %v\n", err.Error())
	}
	return &Service{Name: name, PathPrefix: pathPrefix, DB: db}
}

// ServiceHandler is a function which uses a service reference to produce an http.Handler
type ServiceHandler func(*Service) http.Handler
