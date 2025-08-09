package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"tgbot/lib/e"
	"tgbot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (fp *FetchProcessor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return fp.savePage(chatId, text, username)
	}

	switch text {
	case RndCmd:
		return fp.SendRandom(chatId, username)
	case HelpCmd:
		return fp.SendHelp(chatId)
	case StartCmd:
		return fp.HelloHelp(chatId)
	default:
		return fp.tg.SendMessages(msgUnknownCommand, chatId)
	}
}

func (fp *FetchProcessor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}

	isExists, err := fp.storage.IsExist(page)
	if err != nil {
		return err
	}
	if isExists {
		return fp.tg.SendMessages(msgAlreadyExists, chatID)
	}
	if err := fp.storage.Save(page); err != nil {
		return err
	}
	if err := fp.tg.SendMessages(msgSaved, chatID); err != nil {
		return err
	}
	return nil
}

func (fp *FetchProcessor) SendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't command: rnd page", err) }()

	page, err := fp.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return fp.tg.SendMessages(msgErrSendRandom, chatID)
	} else if errors.Is(err, storage.ErrNoSavedPages) {
		return fp.tg.SendMessages(msgNoPages, chatID)
	}

	if err := fp.tg.SendMessages(page.URL, chatID); err != nil {
		return err
	}

	return fp.storage.Remove(page)
}

func (fp *FetchProcessor) SendHelp(chatID int) (err error) {
	return fp.tg.SendMessages(msgHelp, chatID)
}

func (fp *FetchProcessor) HelloHelp(chatID int) (err error) {
	return fp.tg.SendMessages(msgHello, chatID)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
