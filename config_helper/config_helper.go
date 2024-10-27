package config_helper

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// define struct to access config file
type Config struct {
	Bot struct {
		Token           string `yaml:"token"`
		WebHookURL      string `yaml:"webhook_url"`
		CertFile        string `yaml:"cert_file"`
		KeyFile         string `yaml:"key_file"`
		CronTimePattern string `yaml:"cron_time_pattern"`
	} `yaml:"bot"`
}

func LoadConfig() (*Config, error) {
	var config Config

	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
