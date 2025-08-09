package telegram

import (
	"errors"
	"tgbot/clients/telegram"
	"tgbot/events"
	"tgbot/lib/e"
	"tgbot/storage"
)

type FetchProcessor struct {
	tg      *telegram.Client
	storage storage.Storage
	offset  int
}

type Meta struct {
	ChatId   int
	Username string
}

var ErrUnknownEventType = errors.New("unknown event")
var ErrUnknownMetaType = errors.New("unknown Meta type")

func New(tg *telegram.Client, storage storage.Storage) *FetchProcessor {
	return &FetchProcessor{
		tg:      tg,
		storage: storage,
	}
}

func (fp *FetchProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := fp.tg.Updates(limit, fp.offset)
	if err != nil {
		return nil, e.Wrap("can't get updates", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	events := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		events = append(events, event(update))
	}

	fp.offset = updates[len(updates)-1].ID + 1

	return events, nil
}

func (fp *FetchProcessor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return fp.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (fp *FetchProcessor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	if err := fp.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}
	return nil

}

func meta(ev events.Event) (Meta, error) {
	meta, ok := ev.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("no found meta", ErrUnknownMetaType)
	}
	return meta, nil
}

func event(up telegram.Update) events.Event {
	upType := fetchType(up)

	res := events.Event{
		Type: upType,
		Text: fetchText(up),
	}
	if upType == events.Message {
		res.Meta = Meta{
			ChatId:   up.Message.Chat.Id,
			Username: up.Message.From.Username,
		}
	}

	return res
}

func fetchType(up telegram.Update) events.Type {
	if up.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(up telegram.Update) string {
	if up.Message == nil {
		return ""
	}
	return up.Message.Text
}
