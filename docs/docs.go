package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/cugu/uberfx"
)

//go:embed public
var docs embed.FS

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	fsys, err := fs.Sub(docs, "public")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
}

func main() {
	uberfx.Start(HandleRequest)
}
