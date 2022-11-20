package main

import (
	config "TgMoneyTrackerBot/configuration"
	database "TgMoneyTrackerBot/dbContext"
	messageReceiver "TgMoneyTrackerBot/messageReceiver"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	CreateLogFile()
	db := database.StartDB()
	defer db.Close()
	bot, updatesChannel := StartBot()
	messageReceiver.Start(bot, updatesChannel)
}

func CreateLogFile() {
	logsFolder := "Logs"
	logsPath := path.Join(logsFolder, "Session %d_%d_%d'%dh%dm%ds.txt")
	timeNow := time.Now()
	if _, err := os.Stat(logsFolder); os.IsNotExist(err) {
		os.Mkdir(logsFolder, 0755)
	}

	logsPath = fmt.Sprintf(logsPath,
		timeNow.Day(), timeNow.Month(), timeNow.Year(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())
	file, err := os.OpenFile(logsPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		log.Fatal("error opening log file &s", err)
	}

	log.SetOutput(file)
}

func StartBot() (*tgAPI.BotAPI, tgAPI.UpdatesChannel) {
	bot, err := tgAPI.NewBotAPI(config.GetInfo().BotToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgAPI.NewUpdate(10)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)
	return bot, updates
}
