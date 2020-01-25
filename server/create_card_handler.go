package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/laouji/erinnerung/config"
	"github.com/laouji/erinnerung/slack"
	"github.com/laouji/erinnerung/trello"
	"github.com/laouji/erinnerung/util"
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

	if token := req.FormValue("token"); token != handler.conf.SlackToken {
		handler.writeUnauthorized(w)
		return
	}

	cardName := req.FormValue("text")
	if cardName == "" {
		handler.writeError(w, fmt.Errorf("create card handler requires a card name"))
		return
	}

	today := time.Now()
	dueDate := EstimateDueDate(today)
	card, err := handler.trelloClient.CreateCard(handler.conf.ListId, cardName, dueDate)
	if err != nil {
		handler.writeError(w, fmt.Errorf("create card handler could not create new trello card: %s", err))
		return
	}

	location, err := util.ParseLocation(handler.conf.DisplayTimezone)
	if err != nil {
		handler.writeError(w, fmt.Errorf("create card handler could not load timezone: %s", err))
		return
	}

	payload := &slack.Payload{
		Text:        "OK. Ich erinnere dich an diese Hausaufgabe:",
		UserName:    handler.conf.BotName,
		IconEmoji:   handler.conf.IconEmoji,
		Attachments: slack.GenerateAttachments([]trello.Card{card}, location),
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		handler.writeError(w, fmt.Errorf("create card handler could not marshal payload json: %s", err))
		return
	}

	w.Write(payloadJson)
}

func (handler *CreateCardHandler) writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"text":"an error has occured"}`))
	log.Println(err)
}

func (handler *CreateCardHandler) writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"text":"unauthorized"}`))
}

func EstimateDueDate(dateAssigned time.Time) time.Time {
	durationDay := time.Hour * 24
	dateAssigned = dateAssigned.Truncate(durationDay)

	var daysSinceMonday int
	switch dateAssigned.Weekday() {
	case time.Monday:
		daysSinceMonday = 0
	case time.Sunday:
		daysSinceMonday = 6
	default:
		daysSinceMonday = int(dateAssigned.Weekday()) - 1
	}

	lastMonday := dateAssigned.Add(-time.Duration(daysSinceMonday) * durationDay)
	return lastMonday.Add(7 * durationDay).Add(11 * time.Hour)
}
