package trello

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (client *Client) GetCards(boardId string) (cards []Card, err error) {
	url := fmt.Sprintf("%sboards/%s/cards/open?key=%s&token=%s", client.baseUrl, boardId, client.apiKey, client.apiToken)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return cards, fmt.Errorf("trello client encountered error creating new request: %s", err)
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return cards, fmt.Errorf("trello client encountered error sending request: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return cards, fmt.Errorf("trello client request returned %d status: %s", res.StatusCode, err)
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
