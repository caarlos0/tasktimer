package model

import (
	"encoding/json"
	"log"
	"time"
)

type Task struct {
	ID             uint64    `json:"id"`
	Title          string    `json:"desc"`
	FirstStartedAt time.Time `json:"first_started_at"`
	LastEndedAt    time.Time `json:"last_ended_at"`

	Total     time.Duration   `json:"total_duration"`
	Durations []*TaskDuration `json:"tasks"`
}

type TaskDuration struct {
	StartAt  time.Time     `json:"start"`
	EndAt    time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
}

func (t Task) Bytes() []byte {
	bts, err := json.Marshal(&t)
	if err != nil {
		log.Fatalln(err)
	}
	return bts
}

type ExportedTask struct {
	Title          string    `json:"desc"`
	FirstStartedAt time.Time `json:"first_started_at"`
	LastEndedAt    time.Time `json:"last_ended_at"`

	Total     time.Duration   `json:"total_duration"`
	Durations []*TaskDuration `json:"tasks"`
}

func (t ExportedTask) Bytes() []byte {
	bts, err := json.Marshal(&t)
	if err != nil {
		log.Fatalln(err)
	}
	return bts
}
