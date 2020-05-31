package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/safe-distance/socium-infra/configs/services/proximity/config"
	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/services/proximity/models"
	"github.com/safe-distance/socium-infra/pkg/services/proximity/routes"
)

func main() {
	common.RegisterKafkaClientFlags()
	flag.Parse()
	producer := common.NewObjectLogProducer()
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
