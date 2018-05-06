package main

import (
	"flag"
	"log"

	"github.com/laouji/erinnerung/config"
	"github.com/laouji/erinnerung/reminder"
	"github.com/laouji/erinnerung/slack"
	"github.com/laouji/erinnerung/trello"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.yml", "location of config file")
	flag.Parse()

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("error loading config (%s): %s", configFile, err)
	}

	trelloClient := trello.NewClient(conf.ApiKey, conf.ApiToken)
	slackClient := slack.NewClient(conf.WebHookUri, conf.BotName, conf.IconEmoji)

	reminder := reminder.New(conf, trelloClient, slackClient)
	err = reminder.Remind()
	if err != nil {
		log.Fatalf("error sending reminder %s", err)
	}
}
