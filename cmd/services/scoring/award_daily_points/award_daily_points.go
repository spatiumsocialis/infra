package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron"
	"github.com/spatiumsocialis/infra/configs/services/scoring/config"
	"github.com/spatiumsocialis/infra/pkg/common/kafka"
	"github.com/spatiumsocialis/infra/pkg/services/scoring/models"
)

func main() {
	kafka.RegisterClientFlags()
	flag.Parse()
	p := kafka.NewObjectLogProducer()

	schedule := os.Getenv("SCHEDULE")
	if schedule == "" {
		log.Fatal("SCHEDULE env var not set")
	} else {
		log.Printf("schedule: %v\n", schedule)
	}

	c := cron.New()
	err := c.AddFunc(
		schedule,
		func() {
			dailyAllowance := models.EventScore{
				UID:       config.AllUserID,
				EventID:   1,
				EventType: models.DailyAllowance,
				Timestamp: time.Now(),
				Score:     config.DailyAllowancePoints,
			}
			kafka.LogObject(p, dailyAllowance.UID, dailyAllowance, config.DailyAllowanceTopic)
			fmt.Println("points awarded")
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	c.Run()
}
