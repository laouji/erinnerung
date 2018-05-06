package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/laouji/erinnerung/config"
	"github.com/laouji/erinnerung/trello"
)

type Server struct {
	conf         *config.Data
	trelloClient *trello.Client
}

func New(conf *config.Data, trelloClient *trello.Client) *Server {
	return &Server{
		conf:         conf,
		trelloClient: trelloClient,
	}
}

func (server *Server) Start() {
	mux := http.NewServeMux()
	mux.Handle("/card/create", NewCreateCardHandler(server.conf, server.trelloClient))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", server.conf.Port), mux); err != nil {
		log.Panicf("SERVER SHUTTING DOWN (%s)\n\n", err)
	}
}
