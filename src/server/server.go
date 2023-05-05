package server

import (
	"net/http"
	"xserver/src/config"
)

func AddHandler(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, handler)
}

func Start(config *config.Config) error {
	return http.ListenAndServe(config.Url, nil)
}
