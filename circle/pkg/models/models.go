package models

import (
	"time"

	"github.com/safe-distance/socium-infra/auth"
)

// Circle represents a group of users who aren't making any attempt to social distance from each other (eg families, partners)
type Circle struct {
	ID        string
	Users     []auth.User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
