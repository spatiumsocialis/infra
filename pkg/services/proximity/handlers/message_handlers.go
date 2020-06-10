package handlers

import "github.com/spatiumsocialis/infra/pkg/common/kafka"

// TopicHandlerMap maps topic names to the handlers which handle messages consumed from them
var TopicHandlerMap = map[string]kafka.MessageHandler{
	"user_modified": kafka.SaveUpdatedUserMessageHandler,
}
