package model

import "time"

type Device struct {
	ID               string `json:"id"`               //text eui id for device
	AcceptenceResult bool   `json:"acceptanceResult"` //computed
	Flags            []*Flag
}

type Flag struct {
	Name       string    `json:"name"`
	Number     int       `json:"-"`
	Value      bool      `json:"value"`
	ChangeTime time.Time `json:"changeTimestamp"`
}
