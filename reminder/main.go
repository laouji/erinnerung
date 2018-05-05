package main

import (
	"flag"
	"log"

	"github.com/laouji/erinnerung/config"
	"github.com/laouji/erinnerung/slack"
	"github.com/laouji/erinnerung/trello"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "c", "config.yml", "location of config file")
	flag.Parse()

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("error loading config (%s): %s", configFile, err)
	}

	trelloClient := trello.NewClient(conf.ApiKey, conf.ApiToken)
	cards, err := trelloClient.GetCards(conf.BoardId)
	if err != nil {
		log.Fatalf("error getting trello cards %s", err)
	}

	slackClient := slack.NewClient(conf.WebHookUri, conf.BotName, conf.IconEmoji)
	err = slackClient.Post(cards, conf.DisplayTimezone)
	if err != nil {
		log.Fatalf("error posting to slack %s", err)
	}
}
