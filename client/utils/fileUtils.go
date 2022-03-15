package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const FILENAME = "scanner.exe"

func ParseStopCommand(val string) string {
	return strings.Split(val, "-")[1]
}

func ParseShellCommand(body string) []string {
	return strings.Split(body, " ")
}

func CopyFile(src, dst, filename string) (string, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer source.Close()

	addr := fmt.Sprintf("%s\\%s", dst, filename)
	destination, err := os.Create(addr)
	if err != nil {
		return "", err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return addr, err
}

func GetFileContent(file *os.File) string {
	defer file.Close()
	scanner := bufio.NewScanner(file)
	result := ""
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		result += fmt.Sprintf("%s\n", scanner.Text())
	}

	return result
}

func AppendToFile(filename, text string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s | %s\n", time.Now().String(), text)); err != nil {
		return err
	}
	return nil
}
