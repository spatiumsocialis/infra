package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/location/config"
	"github.com/safe-distance/socium-infra/location/pkg/models"
	"github.com/safe-distance/socium-infra/location/pkg/routes"
)

func main() {
	producer := common.NewNullAsyncProducer()
	common.LoadEnv(false)
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: PORT env variable not set")
	}
	r := common.NewRouter(s, routes.Routes, config.Middleware...)
	http.Handle("/", r)
	fmt.Printf("Listening on port %v...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error listening and serving: %v", err.Error())
	}
}