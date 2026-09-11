package http_mod

import (
	"net/http"
	"time"
)

func NewServer(address string) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           ApplicationRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
