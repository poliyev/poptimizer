package adapters

import (
	"net/http"

	"github.com/WLM1ke/gomoex"
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
