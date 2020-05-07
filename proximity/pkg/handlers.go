package proximity

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddInteraction handles requests for adding new interactions
func AddInteraction(w http.ResponseWriter, r *http.Request) {
	// Get the current user
	user, err := getUser(r)
	if err != nil {
		http.Error(w, "Error retrieving current user: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Decode the interaction from the request body
	var interaction Interaction
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		http.Error(w, "Error decoding interaction from request: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Add the user's UID from the auth token to the interaction
	interaction.UID = user.ID
	db.Create(&interaction)
	json.NewEncoder(w).Encode(&interaction)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

// GetInteractions handles requests to get the current user's interactions
func GetInteractions(w http.ResponseWriter, r *http.Request) {
	// Get the current user
	user, err := getUser(r)
	if err != nil {
		http.Error(w, "Error retrieving current user: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Find the user's interactions and write them to the response
	interactions := make([]Interaction, 0)
	db.Find(&interactions, Interaction{UID: user.ID})
	json.NewEncoder(w).Encode(interactions)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}
