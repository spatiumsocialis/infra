package proximity

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // gorm driver
	"github.com/safe-distance/auth"
	"github.com/safe-distance/initialize"
)

var db *gorm.DB

// Service defines the interface which services must implement
type Service interface {
	Models() []interface{}
	Initialize(envFilename string, overload bool)
	NewRouter() *mux.Router
	PathPrefix() string
}

// InitializeService loads the environment variables and connects to the database and automigrates
func InitializeService(overload bool, envFilenames ...string) {
	initialize.Env(overload, envFilenames...)
	db = initialize.DB(&Interaction{}, &User{})
}

func getUser(r *http.Request) (User, error) {
	// Retrieve the auth token from the request context
	token, err := auth.GetTokenFrom(r.Context())
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

// User is just a gorm model wrapper for the Firebase UID to support circle queries
type User struct {
	ID        string
	CircleID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
