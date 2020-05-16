package models

import (
	"time"

	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/common"
)

// Schema holds the list of models that the DB schema contains
var Schema = common.Schema{
	&Circle{},
	&auth.User{},
}

// Circle represents a group of users who aren't making any attempt to social distance from each other (eg families, partners)
type Circle struct {
	ID        string
	Users     []auth.User
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
