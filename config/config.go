package config

import (
	"os"
)

var (
	Token     string
	BotPrefix string
)

func ReadConfig() error {

	Token = os.Getenv("Token")
	BotPrefix = os.Getenv("BotPrefix")

	return nil
}
