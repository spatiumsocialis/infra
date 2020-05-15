package main

import (
	"flag"

	"github.com/safe-distance/socium-infra/common"

	"github.com/safe-distance/socium-infra/scoring/config"
	"github.com/safe-distance/socium-infra/scoring/pkg/handlers"
	"github.com/safe-distance/socium-infra/scoring/pkg/models"
)

func main() {
	common.RegisterKafkaClientFlags()
	flag.Parse()
	producer := common.NullAsyncProducer{}
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	common.NewConsumer(s, handlers.TopicHandlerMap)
}
