package dto

import (
	"encoding/json"
	"time"
)

type RequestPushEvent struct {
	EventId      string          `json:"event_id"`
	EventOrigin  string          `json:"event_origin"`
	EventDestiny string          `json:"event_destiny"`
	EventType    string          `json:"event_type"`
	Payload      json.RawMessage `json:"payload"`
	Status       string          `json:"status,omitempty"`
	CreatedAt    time.Time       `json:"created_at,omitempty"`
}
