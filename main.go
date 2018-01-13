package main

import (
	"os"

	"github.com/spf13/viper"
	"github.com/pjgg/slack-bot/configuration"
	"github.com/pjgg/slack-bot/commands"
)

func main(){
	configuration.New()
	commands.Execute()
}

func init() {

	viper.SetConfigName("config")
	configPath, exist := os.LookupEnv("CONFIG_PATH")
	if exist {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

}