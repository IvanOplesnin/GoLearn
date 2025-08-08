package main

import (
	"fmt"
	"log"
	"os"
	tg "tgbot/clients/telegram"

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
}
