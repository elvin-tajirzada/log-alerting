package models

import (
	"time"
)

type LokiLogEntry struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream `json:"stream"`
		} `json:"result"`
	} `json:"data"`
}

type Stream struct {
	Caller string    `json:"caller"`
	Dt     string    `json:"dt"`
	IP     string    `json:"ip"`
	Method string    `json:"method"`
	Msg    string    `json:"msg"`
	Path   string    `json:"path"`
	Status string    `json:"status"`
	Timing string    `json:"timing"`
	Job    string    `json:"job"`
	Level  string    `json:"level"`
	Ts     time.Time `json:"ts"`
}
