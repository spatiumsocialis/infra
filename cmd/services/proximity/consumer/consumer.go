package main

import (
	"flag"

	"github.com/spatiumsocialis/infra/configs/services/proximity/config"
	"github.com/spatiumsocialis/infra/pkg/common"
	"github.com/spatiumsocialis/infra/pkg/common/kafka"
	"github.com/spatiumsocialis/infra/pkg/services/proximity/handlers"
	"github.com/spatiumsocialis/infra/pkg/services/proximity/models"
)

func main() {
	kafka.RegisterClientFlags()
	flag.Parse()
	producer := kafka.NullAsyncProducer{}
	s := common.NewService(config.ServiceName, config.ServicePathPrefix, models.Schema, producer, config.ProductionTopic)
	kafka.NewConsumer(s, handlers.TopicHandlerMap)
}
