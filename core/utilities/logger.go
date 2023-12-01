package utilities

import (
	"fmt"
	"log"
	"time"
	"strconv"
	"github.com/AryalKTM/monitor/core/models"
	"github.com/fatih/color"
)

func PrintStatus(service *models.Service, protocol *models.Protocol, err error) (string, error) {
	now := time.Now()
	port := strconv.Itoa(protocol.Port)
	message := ""
	if port == "0" {
		port = "No Port"
	}
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		message = fmt.Sprintf("[%s] [%s] [%s] [%s - %s - %s] An error has occured: %s", red("Error"), now.Format(time.RFC3339), service.Name, protocol.Type, protocol.Server, port, err.Error())
	} else {
		green := color.New(color.FgHiGreen).SprintFunc()
		message = fmt.Sprintf("[%s] [%s] [%s] [%s - %s - %s] Service Seems OK.", green("OK"), now.Format(time.RFC3339), service.Name, protocol.Type, protocol.Server, port)
	}
	log.Println(message)
	return message, err
}

func CreateColorString(str string, clr color.Attribute) string {
	return color.New(clr).SprintFunc()(str)
}