package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spatiumsocialis/infra/configs/services/circle/config"
	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/common/kafka"
	"github.com/spatiumsocialis/infra/pkg/services/circle/models"
	"github.com/spatiumsocialis/infra/pkg/services/circle/routes"
)

func main() {
	kafka.RegisterClientFlags()
	flag.Parse()
	producer := kafka.NewObjectLogProducer()
	common.LoadEnv()
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
