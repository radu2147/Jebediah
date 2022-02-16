package services

import (
	"fmt"
	"main/utils"
)

type SelfDestruct struct {
}

func (s SelfDestruct) Start() {
	err := utils.RemoveFromStartup()
	if err != nil {
		fmt.Println("Error self destructing")
		return
	}
	panic("Goodbye")
}

func (s SelfDestruct) Communicate(_ string) {}
