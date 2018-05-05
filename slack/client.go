package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/laouji/erinnerung/trello"
)

type Client struct {
	httpClient *http.Client
	botName    string
	iconEmoji  string
	webHookUri string
}

func NewClient(webHookUri, botName, iconEmoji string) *Client {
	return &Client{
		httpClient: &http.Client{},
		botName:    botName,
		iconEmoji:  iconEmoji,
		webHookUri: webHookUri,
	}
}

func (client *Client) Post(cards []trello.Card, locationStr string) error {
	if len(cards) == 0 {
		return nil // nothing to do
	}

	location, err := client.parseLocation(locationStr)
	if err != nil {
		return fmt.Errorf("could not load timezone: %s", err)
	}

	attachments := client.generateAttachments(cards, location)
	jsonPayload, err := client.preparePayload(attachments)
	if err != nil {
		return fmt.Errorf("slack client encountered error generating payload: %s", err)
	}

	req, err := http.NewRequest("POST", client.webHookUri, bytes.NewBuffer([]byte(jsonPayload)))
	if err != nil {
		return fmt.Errorf("slack client encountered error creating new request %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.httpClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("slack client encountered error sending request %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("slack client request returned %d status: %s", res.StatusCode, err)
	}

	return nil
}

func (client *Client) preparePayload(attachments []Attachment) (payload []byte, err error) {
	slackPayload := Payload{
		UserName:    client.botName,
		IconEmoji:   client.iconEmoji,
		Text:        "Hallo liebe Leute! Habt ihr schon eure Hausaufgaben fertig gemacht?",
		Attachments: attachments,
	}
	payload, err = json.Marshal(slackPayload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func (client *Client) generateAttachments(cards []trello.Card, location *time.Location) []Attachment {
	attachments := []Attachment{}

	for _, card := range cards {
		attachment := Attachment{
			Title:    card.Name,
			Text:     "_Fälligkeitsdatum_: " + card.Due.In(location).Format("2006-01-02 15:04:05"),
			Fallback: card.Name,
		}

		if card.Desc != "" {
			attachment.TitleLink = card.Desc
		}

		attachments = append(attachments, attachment)
	}

	return attachments
}

func (client *Client) parseLocation(locationStr string) (loc *time.Location, err error) {
	if locationStr == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(locationStr)
}
