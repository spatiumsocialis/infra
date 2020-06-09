package main

import (
	"flag"

	"github.com/spatiumsocialis/infra/configs/services/scoring/config"
	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/common/kafka"
	"github.com/spatiumsocialis/infra/pkg/services/scoring/handlers"
	"github.com/spatiumsocialis/infra/pkg/services/scoring/models"
)

func main() {
	kafka.RegisterClientFlags()
	flag.Parse()
	producer := kafka.NullAsyncProducer{}
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	kafka.NewConsumer(s, handlers.TopicHandlerMap)
}
