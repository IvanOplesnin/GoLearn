package telegram

import "tgbot/clients/telegram"

type FetchProcessor struct {
	tg *telegram.Client
	offset int
}