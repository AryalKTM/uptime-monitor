package protocols

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AryalKTM/uptime-core/models"
)

type Http struct{}

var clientHttp = &http.Client{Timeout: time.Second * 10}

func (httpVar Http) CheckService(Protocol models.Protocol) error {
	url := "http://" + Protocol.Server
	if Protocol.Port == 0 {
		url += ":80"
	} else {
		url += ":" + strconv.Itoa(Protocol.Port)
	}
	resp, err := clientHttp.Get(url)
	if resp != nil {
		_ = resp.Body.Close()
		resp = nil
	}
	return err
}
