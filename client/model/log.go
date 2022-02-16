package model

import "time"

type Log struct {
	Date   time.Time `json:"date"`
	Victim string    `json:"victim"`
}

type TextLog struct {
	Log  Log    `json:"log"`
	Body string `json:"body"`
}

type ImageLog struct {
	Log  Log    `json:"log"`
	Body []byte `json:"body"`
}
