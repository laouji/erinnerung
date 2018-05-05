package trello

import "time"

type Card struct {
	Id     string    `json:"id"`
	Name   string    `json:"name"`
	Desc   string    `json:"desc"`
	Closed bool      `json:"closed"`
	Due    time.Time `json:"due"`
	Url    string    `json:"url"`
}
