package logging

import (
	"time"
)

// LogMessage is the event payload for a log message
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
