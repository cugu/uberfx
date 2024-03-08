package uberfx

import (
	"log"
	"net"
	"net/http"
	"os"
)

func Start(handler http.HandlerFunc) {
	if len(os.Args) != 2 {
		log.Fatal("usage: ./server <address>")
	}

	address := os.Args[1]

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{}

	http.Handle("/", handler)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
