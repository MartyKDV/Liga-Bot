package main

import (
	"fmt"
	"liga-bot/bot"
	"liga-bot/config"
)

func main() {

	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
	}

	bot.Start()
}
