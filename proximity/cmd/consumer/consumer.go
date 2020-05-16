package main

import (
	"flag"

	"github.com/safe-distance/socium-infra/common"
	"github.com/safe-distance/socium-infra/proximity/config"
	"github.com/safe-distance/socium-infra/proximity/pkg/handlers"
	"github.com/safe-distance/socium-infra/proximity/pkg/models"
)

func main() {
	common.RegisterKafkaClientFlags()
	flag.Parse()
	producer := common.NullAsyncProducer{}
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	common.NewConsumer(s, handlers.TopicHandlerMap)
}
