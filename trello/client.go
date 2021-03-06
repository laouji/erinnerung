package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	apiToken   string
	baseUrl    string
}

func NewClient(apiKey, apiToken string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		apiToken:   apiToken,
		baseUrl:    "https://api.trello.com/1/",
	}
}

func (client *Client) CreateCard(listId, taskName string, dueDate time.Time) (card Card, err error) {
	url := fmt.Sprintf("%scards?key=%s&token=%s&idList=%s&name=%s&due=%s",
		client.baseUrl,
		client.apiKey,
		client.apiToken,
		listId,
		url.QueryEscape(taskName),
		dueDate.Format("2006-01-02"),
	)
	res, err := client.sendRequest("POST", url, nil)
	if err != nil {
		return card, err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return card, fmt.Errorf("trello client encountered error reading response body: %s", err)
	}

	err = json.Unmarshal(bodyBytes, &card)
	if err != nil {
		return card, fmt.Errorf("trello client encountered error unmarshalling response body: %s", err)
	}
	return card, nil
}

func (client *Client) GetCards(boardId string) (cards []Card, err error) {
	url := fmt.Sprintf("%sboards/%s/cards/open?key=%s&token=%s", client.baseUrl, boardId, client.apiKey, client.apiToken)
	res, err := client.sendRequest("GET", url, nil)
	if err != nil {
		return cards, err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return cards, fmt.Errorf("trello client encountered error reading response body: %s", err)
	}

	err = json.Unmarshal(bodyBytes, &cards)
	if err != nil {
		return cards, fmt.Errorf("trello client encountered error unmarshalling response body: %s", err)
	}

	return cards, nil
}

func (client *Client) ArchiveOld(cards []Card) error {
	if len(cards) < 1 {
		return nil // nothing to be done
	}

	for _, card := range cards {
		url := fmt.Sprintf("%scards/%s?key=%s&token=%s&closed=true", client.baseUrl, card.Id, client.apiKey, client.apiToken)

		_, err := client.sendRequest("PUT", url, nil)
		if err != nil {
			return fmt.Errorf("trello client encountered error archiving card %s: %s", card.Id, err)
		}
	}

	return nil
}

func (client *Client) sendRequest(method, url string, body io.Reader) (res *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("trello client encountered error creating new request: %s", err)
	}

	res, err = client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("trello client encountered error sending request: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trello client request returned %d status: %s", res.StatusCode, err)
	}

	return res, nil
}
