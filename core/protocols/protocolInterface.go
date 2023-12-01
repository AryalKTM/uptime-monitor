package protocols

import "github.com/AryalKTM/uptime-core/models"

type ProtocolInterface interface {
	CheckService(Protocol models.Protocol) error
}

func RegisterProtocolInterface(registeredInterfaces *map[string]ProtocolInterface, protocolInterfaceName string, protocolInterface ProtocolInterface) {
	(*registeredInterfaces)[protocolInterfaceName] = protocolInterface
}

func DefaultRegisteredProtocolInterfaces() map[string]ProtocolInterface {
	registeredInterfaces := make(map[string]ProtocolInterface, 0)
	registeredInterfaces["https"] = Https{}
	registeredInterfaces["http"] = Http{}
	registeredInterfaces["ftp"] = FTP{}
	return registeredInterfaces
}

func IsRegistered(registeredInterfaces *map[string]ProtocolInterface, funcName string) bool {
	interf := (*registeredInterfaces)[funcName]

	if interf != nil {
		return true
	} else {
		return false
	}
}
