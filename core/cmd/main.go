package main

import (
	"os"
	"github.com/AryalKTM/uptime-core/responsehandlers"
	"github.com/AryalKTM/uptime-core/protocols"
	"github.com/AryalKTM/uptime-core/core"
)

func main() {
	regProtocoInterfaces, regResponseHandlerInterfaces := core.Initialize(os.Args[1])

	protocols.RegisterProtocolInterface(&regProtocoInterfaces, "ftp", protocols.FTP{})

	responsehandlers.RegisterResponseHandlerInterface(&regResponseHandlerInterfaces, "consoleWithMemory", responsehandlers.ConsoleHandlerWithMemory{})

	core.Start(regProtocoInterfaces, regResponseHandlerInterfaces)
}