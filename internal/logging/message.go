package logging

import (
	"time"
)

type LogMessage struct {
	ID          string
	Time        time.Time
	Level       string
	Persist     bool
	PersistTime time.Duration
	Attributes  []Attr
	Message     string `json:"msg"`
}

type Attr struct {
	Key   string
	Value string
}
