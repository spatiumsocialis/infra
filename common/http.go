package common

import (
	"log"
	"net/http"
)

// ThrowError logs an error and writes it to the response
func ThrowError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
