package common

import "github.com/jinzhu/gorm"

// Service represents an HTTP service
type Service struct {
	Name       string
	DB         *gorm.DB
	PathPrefix string
}

// NewService returns a new service constructed with the given arguments
func NewService(name string, db *gorm.DB, pathPrefix string) *Service {
	return &Service{Name: name, DB: db, PathPrefix: pathPrefix}
}
