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

// Schema is a type alias for an empty-interface slice, for storing schema model references
type Schema []interface{}

// NewService returns a new service constructed with the given arguments
func NewService(name string, pathPrefix string, schema Schema) *Service {
	db, err := NewDB(schema...)
	if err != nil {
		log.Fatalf("Error initializing DB: %v\n", err.Error())
	}
	return &Service{Name: name, PathPrefix: pathPrefix, DB: db}
}

// ServiceHandler is a function which uses a service reference to produce an http.Handler
type ServiceHandler func(*Service) http.Handler
