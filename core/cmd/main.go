package main

import (
	"os"
)

func main() {
	regProtocoInterfaces, regResponseHandlerInterfaces := core.Initialize(os.Args[1])

	protocols.RegsiterProtocolInterface(&regProtocoInterfaces, "ftp", protocols.FTP{})

	responsehandlers.RegisterResponseHandlerInterface(&regResponseHandlerInterfaces, "consoleWithMemory", responsehandlers.ConsoleHandlerWithMemory{})

	core.Start(regProtocoInterfaces, regResponseHandlerInterfaces)
}