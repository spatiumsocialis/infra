package common

import (
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
func NewService(name string, pathPrefix string, db *gorm.DB) *Service {
	return &Service{Name: name, PathPrefix: pathPrefix, DB: db}
}

// ServiceHandler is a function which uses a service reference to produce an http.Handler
type ServiceHandler func(*Service) http.Handler
