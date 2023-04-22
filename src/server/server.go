package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type ServerOptions struct {
	Url string
}

type Server struct {
	options *ServerOptions
	router  *mux.Router
}

func Create(options *ServerOptions) *Server {
	return &Server{
		options: options,
		router:  mux.NewRouter(),
	}
}

func (xserver *Server) AddHandler(path string, handler http.HandlerFunc, methood string) {
	xserver.router.HandleFunc(path, handler)
}

func (xserver *Server) Start() error {
	return http.ListenAndServe(xserver.options.Url, xserver.router)
}
