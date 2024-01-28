//go:build !wasip1

package wasi

import (
	"crypto/tls"
	"net/http"

	_ "github.com/stealthrocket/net/http"
)

func Client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
