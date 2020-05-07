package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	proximity "github.com/safe-distance/proximity/pkg"
)

func main() {
	proximity.InitializeService(false)
	port := os.Getenv("PORT")
	// servicePrefix := os.Getenv("SERVICE_PREFIX")
	if port == "" {
		log.Fatal("Error: PORT env variable not set")
	}
	r := proximity.NewRouter()
	http.Handle("/", r)
	fmt.Printf("service prefix: %v\n", proximity.ServicePrefix)
	fmt.Println("Listening...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error listening and serving: %v", err.Error())
	}
}
