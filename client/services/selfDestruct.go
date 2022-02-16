package services

import "main/utils"

type SelfDestruct struct {
}

func (s SelfDestruct) Start() {
	err := utils.RemoveFromStartup()
	if err != nil {
		return
	}
	panic("Goodbye")
}

func (s SelfDestruct) Communicate(_ string) {}
