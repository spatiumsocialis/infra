package main

import (
	"flag"

	"github.com/safe-distance/socium-infra/configs/services/scoring/config"
	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/services/scoring/handlers"
	"github.com/safe-distance/socium-infra/pkg/services/scoring/models"
)

func main() {
	common.RegisterKafkaClientFlags()
	flag.Parse()
	producer := common.NullAsyncProducer{}
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	common.NewConsumer(s, handlers.TopicHandlerMap)
}
