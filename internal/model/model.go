package model

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID        uint64        `json:"id"`
	Title     string        `json:"desc"`
	StartAt   time.Time     `json:"start"`
	EndAt     time.Time     `json:"end"`
	PausedAt  time.Time     `json:"paused_at,omitempty"`
	PausedFor time.Duration `json:"paused_for,omitempty"`
}

func (t Task) Bytes() ([]byte, error) {
	return json.Marshal(&t)
}

type ExportedTask struct {
	Title   string    `json:"desc"`
	StartAt time.Time `json:"start"`
	EndAt   time.Time `json:"end"`
}

func (t ExportedTask) Bytes() ([]byte, error) {
	return json.Marshal(&t)
}
