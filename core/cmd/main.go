package main

import (
	"os"
	"github.com/AryalKTM/monitor/core/responsehandlers"
	"github.com/AryalKTM/monitor/core/protocols"
	"github.com/AryalKTM/monitor/core/core"
)

func main() {
	regProtocoInterfaces, regResponseHandlerInterfaces := core.Initialize(os.Args[0])

	protocols.RegisterProtocolInterface(&regProtocoInterfaces, "ftp", protocols.FTP{})

	responsehandlers.RegisterResponseHandlerInterface(&regResponseHandlerInterfaces, "consoleWithMemory", responsehandlers.ConsoleHandlerWithMemory{})

	core.Start(regProtocoInterfaces, regResponseHandlerInterfaces)
}