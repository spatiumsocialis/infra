package handlers

import (
	"github.com/safe-distance/socium-infra/common"
)

// TopicHandlerMap maps topic names to the handlers which handle messages consumed from them
var TopicHandlerMap = map[string]common.MessageHandler{
	"user_modified": common.SaveUpdatedUserMessageHandler,
}
