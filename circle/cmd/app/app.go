package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/safe-distance/socium-infra/auth"
	"github.com/safe-distance/socium-infra/circle/config"
	"github.com/safe-distance/socium-infra/circle/pkg/models"
	"github.com/safe-distance/socium-infra/circle/pkg/routes"
	"github.com/safe-distance/socium-infra/common"
)

func main() {
	common.LoadEnv(false)
	db, err := common.NewDB(&models.Circle{}, &auth.User{})
	if err != nil {
		log.Fatalf("Error initializing DB: %v\n", err.Error())
	}
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, db)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: PORT env variable not set")
	}
	r := routes.NewRouter(s)
	http.Handle("/", r)
	fmt.Println("Listening...")
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error listening and serving: %v", err.Error())
	}
}
