package commands

import (
	"github.com/spf13/cobra"
	"github.com/pjgg/slack-bot/connectors"
	"github.com/pjgg/slack-bot/configuration"
)

var (
	listenerCmd = &cobra.Command{
		Use: "slack-listener",
		Run: slackListenerHandler,
	}
)

func init() {
	BaseCmd.AddCommand(listenerCmd)
}

func slackListenerHandler(cmd *cobra.Command, args []string) {
	slackConnector := connectors.Instance(configuration.ConfigurationManagerInstance.SlackToken)
	slackConnector.SlackBotListener()
}