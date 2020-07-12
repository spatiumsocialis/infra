package auth

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

type (
	// User is just a gorm model wrapper for the Firebase UID to support circle queries
	User struct {
		ID        string     `json:"id"`
		CircleID  string     `json:"circleId"`
		CreatedAt time.Time  `json:"-"`
		UpdatedAt time.Time  `json:"-"`
		DeletedAt *time.Time `json:"-"`
	}
)

// GetUser retrieves the UID from a request's token and returns the user model
func GetUser(db *gorm.DB, r *http.Request) (User, error) {
	// Retrieve the auth token from the request context
	token, err := GetToken(r.Context())
	if err != nil {
		return User{}, err
	}
	// Retrieve the client UID from the token
	currentUID := token.UID

	// Fetch the current user
	var user User
	if result := db.FirstOrCreate(&user, User{ID: currentUID}); result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}
