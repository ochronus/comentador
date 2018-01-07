package main

import (
	"time"
)

type ConfigSpec struct {
	Debug      bool
	Port       int
	User       string
	AkismetKey string
	Users      []string
	Rate       float32
	Timeout    time.Duration
	ColorCodes map[string]int
}
