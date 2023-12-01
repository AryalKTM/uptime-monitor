package responsehandler

import (
	"github.com/AryalKTM/monitor/core/models"
	"github.com/AryalKTM/monitor/core/models/response"
	"github.com/AryalKTM/monitor/monitor/models/data"
	"github.com/fatih/color"
	"strconv"
	"time"
	"fmt"
	"log"
)

type WebServerRespHandler struct {
	ServiceProtocol map[string]response.ResponseType
	OutputChannel   chan response.Response
}

func (handler WebServerRespHandler) Handle(configuration *models.Config, channel *chan response.Response) {
	handler.loadServices(configuration)
	handler.ServiceProtocol = make(map[string]response.ResponseType)
	go func() {
		for {
			select {
			case resp := <-*channel:
				handler.writeResponse(resp)
			}
		}
	}()
}

func (handler WebServerRespHandler) writeResponse(response response.Response) {
	select {
	case handler.OutputChannel <- response:
		return
	}
}

func (handler WebServerRespHandler) loadServices(configuration *models.Config) {
	for _, service := range configuration.Services {
		for _, protocol := range service.Protocols {
			resp := response.Response{
				ServiceName: service.Name,
				Protocol:    protocol,
				Error:       data.NewCustomErr("<<pending>>"),
			}
			handler.writeResponse(resp)
		}
	}
}

func printIfChange(resp response.Response, c *WebServerRespHandler) {
	respType := c.ServiceProtocol[resp.ServiceName+resp.Protocol.Type]
	if respType != resp.ResponseType {
		c.ServiceProtocol[resp.ServiceName+resp.Protocol.Type] = resp.ResponseType
		now := time.Now()
		port := strconv.Itoa(resp.Protocol.Port)
		message := ""
		if port == "0" {
			port = "No Port"
		}
		if resp.ResponseType == response.Error {
			red := color.New(color.FgRed).SprintFunc()
			message = fmt.Sprintf("[%s] [%s] [%s] [%s - %s - %s] An Error has Occured: %s", red("Error"), now.Format(time.RFC3339), resp.ServiceName, resp.Protocol.Type, resp.Protocol.Server, port, resp.Error.Error())
		} else {
			green := color.New(color.FgHiGreen).SprintFunc()
			message = fmt.Sprintf("[%s] [%s] [%s] [%s - %s - %s] Service Seems OK.", green("OK"), now.Format(time.RFC3339), resp.ServiceName, resp.Protocol.Type, resp.Protocol.Server, port)
		}
		log.Println(message)
	}
}
