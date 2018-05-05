package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Data struct {
	BotName         string `yaml:"bot_name"`
	IconEmoji       string `yaml:"icon_emoji"`
	DisplayTimezone string `yaml:"display_timezone"`
	WebHookUri      string `yaml:"web_hook_uri"`
	ApiKey          string `yaml:"api_key"`
	ApiToken        string `yaml:"api_token"`
	BoardId         string `yaml:"board_id"`
}

func LoadConfig(filepath string) (data *Data, err error) {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read conf file: %s", err)
	}

	data = &Data{}
	if err := yaml.Unmarshal(buf, data); err != nil {
		return nil, fmt.Errorf("could not parse yaml: %s", err)
	}

	return data, nil
}
