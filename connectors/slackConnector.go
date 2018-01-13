package connectors

import (
	"os"
	"fmt"
	"sync"
	"strings"
	"bufio"
	"os/exec"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	
)

// SlackConnector ...struct contains a slack client pointer and a real time protocol pointer
type SlackConnector struct { 
	Client *slack.Client
	RealTimeMsgProtocol *slack.RTM
}

type SlackConnectorBehavior interface {
	SlackBotListener()
}

var onceSlack sync.Once
// SlackConnectorInstance is a slack client singleton instance 
var SlackConnectorInstance SlackConnector

// New ...given a slack bot token, returns slack connector single instance.
func Instance(token string) *SlackConnector{
	onceSlack.Do(func() {
		SlackConnectorInstance.Client = slack.New(token)
		SlackConnectorInstance.RealTimeMsgProtocol =  SlackConnectorInstance.Client.NewRTM()
		
		go SlackConnectorInstance.RealTimeMsgProtocol.ManageConnection()
	})

	return &SlackConnectorInstance
}

func (sc *SlackConnector) SlackBotListener(){
	for {
		select {
		case msg := <-sc.RealTimeMsgProtocol.IncomingEvents:
			logrus.Info("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				logrus.Info("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				info := sc.RealTimeMsgProtocol.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)
				commandArgs := strings.Replace(ev.Text,prefix,"",-1)
				logrus.Info("Command: \n", commandArgs)
				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
				   pwd, _ := os.Getwd()
				   cmd := exec.Command(pwd+ "/main",commandArgs)
				   cmdReader, err := cmd.StdoutPipe()
				   if err != nil{
				   	logrus.Error("Error creating StdoutPipe: %v", err)
				   }
				   
					scanner := bufio.NewScanner(cmdReader)
					go func() {
						for scanner.Scan() {
							sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(scanner.Text(), ev.Channel))
						}
					}()

					if err := cmd.Start(); err != nil {
						logrus.Error(err)
						}


					if err := cmd.Wait(); err != nil {
						logrus.Error(err)
					}
				   
				}

			case *slack.RTMError:
				logrus.Error("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				logrus.Error("Invalid credentials")

			default:
				// do nothing
			}
		}
	}
}