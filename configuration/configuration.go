package configuration

import (
	"sync"

	"github.com/spf13/viper"
)

type ConfigurationManager struct {
	SlackToken string
}

var once sync.Once
var ConfigurationManagerInstance *ConfigurationManager

func New() *ConfigurationManager {

	once.Do(func() {
		viper.BindEnv("slack.token", "SLACK_TOKEN")
		slackToken := viper.GetString("slack.token")
		ConfigurationManagerInstance = &ConfigurationManager{
			SlackToken: slackToken,
		}
	})

	return ConfigurationManagerInstance
}