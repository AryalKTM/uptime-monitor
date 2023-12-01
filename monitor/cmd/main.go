package main

import (
	"os"
	"log"

	"github.com/AryalKTM/monitor/core/core"
	"github.com/AryalKTM/monitor/core/responsehandlers"
	"github.com/AryalKTM/monitor/core/utilities"
	"github.com/AryalKTM/monitor/monitor/system"
	"github.com/AryalKTM/monitor/monitor/webserver"
	"github.com/AryalKTM/monitor/monitor/webserver/responsehandler"
	"github.com/fatih/color"
)

func main() {
		errEnv := os.Setenv("GHW_DISABLE_WARNINGS", "1")
		if errEnv != nil {
			log.Printf("[%s] %s", utilities.CreateColorString("Warning",color.FgHiYellow), errEnv)
		}
		system.PrintSystemInfo()
		ws := webserver.NewWebServer("config/serverconfig.yml")
		ws.Start()

		regProtocolInterfaces, regResponseHandlerInterfaces := core.Initialize("config/serverconfig.yml")

		responsehandlers.RegisterResponseHandlerInterface(&regResponseHandlerInterfaces, "webServerHandler", responsehandler.WebServerRespHandler{OutputChannel: ws.InputChannel})
		responsehandlers.RegisterResponseHandlerInterface(&regResponseHandlerInterfaces, "consoleMemory", responsehandlers.ConsoleHandlerWithMemory{})

		core.Start(regProtocolInterfaces, regResponseHandlerInterfaces)
}