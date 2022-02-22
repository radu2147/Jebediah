package models

import "time"

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
