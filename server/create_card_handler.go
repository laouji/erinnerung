package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/laouji/erinnerung/config"
	"github.com/laouji/erinnerung/slack"
	"github.com/laouji/erinnerung/trello"
)

type CreateCardHandler struct {
	conf         *config.Data
	trelloClient *trello.Client
}

func NewCreateCardHandler(conf *config.Data, trelloClient *trello.Client) *CreateCardHandler {
	return &CreateCardHandler{
		conf:         conf,
		trelloClient: trelloClient,
	}
}

func (handler *CreateCardHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	payload := &slack.Payload{
		Text:      "test response",
		UserName:  handler.conf.BotName,
		IconEmoji: handler.conf.IconEmoji,
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, `{"text":"an error has occured"}`, http.StatusInternalServerError)
		log.Println("could not marshal payload json", err)
		return
	}

	w.Write(payloadJson)
}
