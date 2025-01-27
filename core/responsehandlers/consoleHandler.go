package responsehandlers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/AryalKTM/monitor/core/models"
	"github.com/AryalKTM/monitor/core/models/response"
	"github.com/fatih/color"
)

type ConsoleHandler struct{}

func (handler ConsoleHandler) Handle(configuration *models.Config, channel *chan response.Response) {
	for {
		resp := <-*channel
		print(resp)
	}
}

func print(resp response.Response) {
	now := time.Now()
	port := strconv.Itoa(resp.Protocol.Port)
	message := ""
	if port == "0" {
		port = "No Port"
	}
	if resp.ResponseType == response.Error {
		red := color.New(color.FgRed).SprintFunc()
		message = fmt.Sprintf("[%s] [%s] [%s] [%s - %s - %s] An error as accured: %s", red("ERRO"), now.Format(time.RFC3339), resp.ServiceName, resp.Protocol.Type, resp.Protocol.Server, port, resp.Error.Error())
	} else {
		green := color.New(color.FgHiGreen).SprintFunc()
		message = fmt.Sprintf("[ %s ] [%s] [%s] [%s - %s - %s] Service seems OK.", green("OK"), now.Format(time.RFC3339), resp.ServiceName, resp.Protocol.Type, resp.Protocol.Server, port)
	}
	log.Println(message)
}
