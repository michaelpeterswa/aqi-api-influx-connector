package requests

import (
	"net/http"
	"time"
)

type RequestClient struct {
	client *http.Client
}

func NewRequestClient() *RequestClient {
	return &RequestClient{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}
