package main

import (
	"log"
	"os"
	tg "tgbot/clients/telegram"
	"tgbot/consumer/event_consumer"
	processTg "tgbot/events/telegram"
	"tgbot/storage/files"

	"github.com/joho/godotenv"
)

var TOKEN string
var HOST string

func init() {
	godotenv.Load()
	TOKEN = os.Getenv("TOKEN")
	HOST = os.Getenv("HOST")
}

func main() {
	if TOKEN == "" {
		log.Fatal("No TOKEN in environments")
	}

	tgClient := tg.New(HOST, TOKEN)

	storage := files.New("/data_bot")

	eventsProcessor := processTg.New(tgClient, storage)

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, 10)

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}
