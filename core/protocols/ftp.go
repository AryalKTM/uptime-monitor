package protocols

import (
	"github.com/AryalKTM/uptime-core/models"
)

type FTP struct {
}
	func (ftp FTP) CheckService(Protocol models.Protocol) error {
		return nil
}