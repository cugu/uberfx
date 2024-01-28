//go:build wasip1

package wasi

import (
	"crypto/tls"
	"net/http"

	_ "github.com/stealthrocket/net/http"
	"github.com/stealthrocket/net/wasip1"
)

func Client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext:     wasip1.DialContext,
		},
	}
}
