package reminder

import (
	"fmt"
	"time"

	"github.com/laouji/erinnerung/config"
	"github.com/laouji/erinnerung/slack"
	"github.com/laouji/erinnerung/trello"
)

type Reminder struct {
	conf         *config.Data
	trelloClient *trello.Client
	slackClient  *slack.Client
}

func New(conf *config.Data, trelloClient *trello.Client, slackClient *slack.Client) *Reminder {
	return &Reminder{
		conf:         conf,
		trelloClient: trelloClient,
		slackClient:  slackClient,
	}
}

func (r *Reminder) Remind() error {
	cards, err := r.trelloClient.GetCards(r.conf.BoardId)
	if err != nil {
		return fmt.Errorf("failed to get trello cards: %s", err)
	}

	currentCards, pastDueCards := r.filter(cards)

	err = r.slackClient.Post(currentCards, r.conf.DisplayTimezone)
	if err != nil {
		return fmt.Errorf("failed to post to slack: %s", err)
	}

	err = r.trelloClient.ArchiveOld(pastDueCards)
	if err != nil {
		return fmt.Errorf("failed to archive old: %s", err)
	}
	return nil
}

func (r *Reminder) filter(cards []trello.Card) (current, pastDue []trello.Card) {
	now := time.Now()
	current = []trello.Card{}
	pastDue = []trello.Card{}

	for _, card := range cards {
		if card.Due.Before(now) {
			pastDue = append(pastDue, card)
			continue
		}
		current = append(current, card)
	}
	return current, pastDue
}
