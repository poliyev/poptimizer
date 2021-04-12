package adapters

import (
	"github.com/WLM1ke/gomoex"
	"net/http"
)

// NewISSClient - создает клиент для ISS с ограничением на количество соединений.
func NewISSClient(maxCons int) *gomoex.ISSClient {
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: maxCons,
		},
	}
	return gomoex.NewISSClient(client)
}
