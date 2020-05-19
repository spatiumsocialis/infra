package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron"
	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/scoring/config"
	"github.com/safe-distance/socium-infra/scoring/pkg/models"
)

func main() {
	common.RegisterKafkaClientFlags()
	flag.Parse()
	p := common.NewObjectLogProducer()
	dailyAllowance := models.EventScore{
		UID:       config.AllUserID,
		EventID:   1,
		EventType: models.DailyAllowance,
		Timestamp: time.Now(),
		Score:     config.DailyAllowancePoints,
	}
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
			common.LogObject(p, dailyAllowance.UID, dailyAllowance, config.DailyAllowanceTopic)
			fmt.Println("points awarded")
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	c.Run()
}
