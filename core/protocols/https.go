package protocols

import (
	"net/http"
	"time"
	"strconv"

	"github.com/AryalKTM/monitor/core/models"
)

type Https struct {}

var clientHTTPS = &http.Client{Timeout: time.Second * 10}

func (https Https) CheckService(Protocol models.Protocol) error {
	url := "https://" + Protocol.Server
	if Protocol.Port == 0 {
		url += ":443"
	} else {
		url += ":" + strconv.Itoa(Protocol.Port)
	}
	resp, err := clientHTTPS.Get(url)
	if resp != nil {
		_ = resp.Body.Close()
		resp = nil
	}
	return err
}