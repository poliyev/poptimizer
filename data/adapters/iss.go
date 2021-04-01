package adapters

import (
	"github.com/WLM1ke/gomoex"
	"net/http"
	"time"
)

func NewISSClient() *gomoex.ISSClient {
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: 20,
		},
		Timeout: 30 * time.Second,
	}
	return gomoex.NewISSClient(client)
}
