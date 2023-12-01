package responsehandlers

import (
	"github.com/AryalKTM/monitor/core/models"
	"github.com/AryalKTM/monitor/core/models/response"
)

type ResponseHandlerInterface interface {
	Handle(*models.Config, *chan response.Response)
}

func RegisterResponseHandlerInterface(registeredResponseHandlerInterfaces *map[string]ResponseHandlerInterface, responseInterfaceName string, responseHandlerInterface ResponseHandlerInterface) {
	(*registeredResponseHandlerInterfaces)[responseInterfaceName] = responseHandlerInterface
}

func DefaultRegisteredResponseHandlers() map[string]ResponseHandlerInterface {
	registeredHandlers := make(map[string]ResponseHandlerInterface, 0)
	return registeredHandlers
}