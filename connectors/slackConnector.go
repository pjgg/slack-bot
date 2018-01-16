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
				// Like taxi driver movie!. 'are you talking to me?'
				if sc.isTalkingToMe(ev) {
				   cmd := exec.Command(os.Args[0],sc.getUserExecCommand(ev))
				   sc.replyCommandStd(cmd, ev)
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

func (sc *SlackConnector) getEventUserIDTag(info *slack.Info) string {
	return fmt.Sprintf("<@%s> ", info.User.ID)
}

func (sc *SlackConnector) getUserExecCommand(ev *slack.MessageEvent) string {
	info := sc.RealTimeMsgProtocol.GetInfo()
	return strings.Replace(ev.Text,sc.getEventUserIDTag(info),"",-1)
}

func (sc *SlackConnector) isTalkingToMe(ev *slack.MessageEvent) bool {
	info := sc.RealTimeMsgProtocol.GetInfo()
	return ev.User != info.User.ID && strings.HasPrefix(ev.Text, sc.getEventUserIDTag(info))
}

func (sc *SlackConnector) replyCommandStd(cmd *exec.Cmd, ev *slack.MessageEvent){
	cmdReader, err := cmd.StdoutPipe()
	if err != nil{
		logrus.Error("Error creating StdoutPipe: %v", err)
		sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(err.Error(), ev.Channel))
	}
	
	 scanner := bufio.NewScanner(cmdReader)
	 go func() {
		 for scanner.Scan() {
			 sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(scanner.Text(), ev.Channel))
		 }
	 }()

	 if err := cmd.Start(); err != nil {
		 sc.RealTimeMsgProtocol.SendMessage(sc.RealTimeMsgProtocol.NewOutgoingMessage(err.Error(), ev.Channel))
		 logrus.Error(err)
		 }


	 if err := cmd.Wait(); err != nil {
		 logrus.Error(err)
	 }
}