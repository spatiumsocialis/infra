package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/proximity/config"
	"github.com/safe-distance/socium-infra/proximity/pkg/routes"
)

func main() {
	common.LoadEnv(false)
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, config.Models...)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: PORT env variable not set")
	}
	r := common.NewRouter(s, routes.Routes, config.Middleware...)
	http.Handle("/", r)
	fmt.Println("Listening...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error listening and serving: %v", err.Error())
	}
}
