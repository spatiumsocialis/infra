package main

import (
	"flag"

	"github.com/safe-distance/socium-infra/configs/services/proximity/config"
	"github.com/safe-distance/socium-infra/pkg/common"
	"github.com/safe-distance/socium-infra/pkg/common/kafka"
	"github.com/safe-distance/socium-infra/pkg/services/proximity/handlers"
	"github.com/safe-distance/socium-infra/pkg/services/proximity/models"
)

func main() {
	kafka.RegisterClientFlags()
	flag.Parse()
	producer := kafka.NullAsyncProducer{}
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	kafka.NewConsumer(s, handlers.TopicHandlerMap)
}
