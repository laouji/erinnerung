package main

import (
	"flag"
	"log"

	"github.com/laouji/erinnerung/config"
	"github.com/laouji/erinnerung/server"
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

	server := server.New(conf, trelloClient)
	server.Start()
}
