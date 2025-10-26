package domain

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventId      string
	EventOrigin  string
	EventDestiny string
	EventType    string
	Payload      json.RawMessage
	Status       string
	CreatedAt    time.Time
}
