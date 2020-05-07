package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	circle "github.com/safe-distance/circle/pkg/circle"
)

func main() {
	circle.InitializeService(false)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: PORT env variable not set")
	}
	r := circle.NewRouter()
	http.Handle("/", r)
	fmt.Println("Listening...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error listening and serving: %v", err.Error())
	}
}
