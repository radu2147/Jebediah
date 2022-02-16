package main

import (
	"fmt"
	"os"
	"time"
)

func LogEntry(file *os.File, log *TextLog, raw bool) error {
	if raw {
		_, err := fmt.Fprintf(file, "%s", log.Body)
		if err != nil {
			print(err)
			return err
		}
	} else {
		_, err := fmt.Fprintf(file, "%s | %s\n", time.Now().String(), log.Body)
		if err != nil {
			print(err)
			return err
		}
	}
	return nil
}
