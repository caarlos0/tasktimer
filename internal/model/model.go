package model

import (
	"encoding/json"
	"log"
	"time"
)

type Task struct {
	ID      uint64    `json:"id"`
	Title   string    `json:"desc"`
	StartAt time.Time `json:"start"`
	EndAt   time.Time `json:"end"`
}

func (t Task) Bytes() []byte {
	bts, err := json.Marshal(&t)
	if err != nil {
		log.Fatalln(err)
	}
	return bts
}

type ExportedTask struct {
	Title   string    `json:"desc"`
	StartAt time.Time `json:"start"`
	EndAt   time.Time `json:"end"`
}

func (t ExportedTask) Bytes() []byte {
	bts, err := json.Marshal(&t)
	if err != nil {
		log.Fatalln(err)
	}
	return bts
}
