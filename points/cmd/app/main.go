package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/safe-distance/auth"
	"github.com/safe-distance/points"
)

func main() {
	points.InitializeService()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: PORT env variable not set")
	}
	r := mux.NewRouter()
	r.HandleFunc("/scores", points.Handler)
	r.Use(auth.Middleware)
	http.Handle("/", r)
	fmt.Printf("Listening on port %v...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error listening and serving: %v", err.Error())
	}
}
