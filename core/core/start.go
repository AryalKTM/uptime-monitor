package core

import (
	"log"
	"strconv"
	"time"

	"github.com/AryalKTM/monitor/core/models"
	"github.com/AryalKTM/monitor/core/models/response"
	"github.com/AryalKTM/monitor/core/process"
	"github.com/AryalKTM/monitor/core/protocols"
	"github.com/AryalKTM/monitor/core/responsehandlers"
	"github.com/AryalKTM/monitor/core/utilities"
	"github.com/fatih/color"
)

var _configPath string

func Initialize(configPath string) (map[string]protocols.ProtocolInterface, map[string]responsehandlers.ResponseHandlerInterface) {
	_configPath = configPath
	return protocols.DefaultRegisteredProtocolInterfaces(), responsehandlers.DefaultRegisteredResponseHandlers()
}

func Start(registeredProtocolInterfaces map[string]protocols.ProtocolInterface, registeredResponseHandlerInterfaces map[string]responsehandlers.ResponseHandlerInterface) {
	conf, err := models.ConfigFromFile(_configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[%s] Configuration Loaded!!!", utilities.CreateColorString("Load", color.FgHiBlue))
	processesChannel := make(chan string)
	responseChannel := make(chan response.Response)

	startResponseBrodadcaster(conf, &responseChannel, &registeredResponseHandlerInterfaces)
	log.Printf("[%s] Starting Handlers Complete!", utilities.CreateColorString("Complete", color.FgHiBlue))

	if err != nil {
		log.Fatal(err.Error())
	}

	for i := range conf.Services {
		service := conf.Services[i]
		for i2 := range service.Protocols {
			protocol := service.Protocols[i2]
			protocolInterface := registeredProtocolInterfaces[protocol.Type]
			if protocols.IsRegistered(&registeredProtocolInterfaces, protocol.Type) {
				proc := process.NewProcess(func() {
					if err = protocolInterface.CheckService(protocol); err == nil {
						responseChannel <- response.Response{
							ServiceName:  service.Name,
							Protocol:     protocol,
							ResponseType: response.Success,
							Error:        nil}
					} else {
						responseChannel <- response.Response{
							ServiceName:  service.Name,
							Protocol:     protocol,
							ResponseType: response.Error,
							Error:        err}
					}
				}, processesChannel)
				process.ScheduleProcess(proc, protocol.Interval)
			} else {
				red := color.New(color.FgRed).SprintFunc()
				log.Printf("[%s] [%s] [%s] [%s - %s - %s] An Error Has Occured: %s", red("Error"), time.Now().Format(time.RFC3339), service.Name, protocol.Type, protocol.Server, strconv.Itoa(protocol.Port), "Protocol Interface"+protocol.Type+"Not Registered")
			}

	}
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
}

func startResponseBrodadcaster(configuration *models.Config, responseChannel *chan response.Response, responseHandlers *map[string]responsehandlers.ResponseHandlerInterface) {
	chanArray := make([]*chan response.Response, 0)

	for _, handler := range *responseHandlers {
		channel := make(chan response.Response)
		chanArray = append(chanArray, &channel)

		handler := handler
		go func() {
			handler.Handle(configuration, &channel)
		}()
	}

	go func() {
		for {
			resp := <-*responseChannel
			for _, channel := range chanArray {
				*channel <- resp
			}
		}
	}()
}
