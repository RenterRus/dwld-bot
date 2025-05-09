package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type buttons struct {
	ID   string
	Text string
}

const (
	ActualState     = "ActualState"
	CleanHistory    = "CleanHistory"
	RemoveFromQueue = "RemoveFromQueue"
	ViewQueue       = "ViewQueue"
	LinksForUtil    = "LinksForUtil"
)

func (b *Bot) initKeyboard() map[string]buttons {
	buttonsMap := make(map[string]buttons)

	buttonsMap[ActualState] = buttons{
		ID:   ActualState,
		Text: "Текущее состояние",
	}
	buttonsMap[CleanHistory] = buttons{
		ID:   CleanHistory,
		Text: "Очистить историю",
	}
	buttonsMap[RemoveFromQueue] = buttons{
		ID:   RemoveFromQueue,
		Text: "Удалить из очереди",
	}
	buttonsMap[ViewQueue] = buttons{
		ID:   ViewQueue,
		Text: "Показать очередь",
	}
	buttonsMap[LinksForUtil] = buttons{
		ID:   LinksForUtil,
		Text: "Список ссылок в работе (для утилиты)",
	}

	return buttonsMap
}

func (k *Bot) keyboardDefault() tgbotapi.InlineKeyboardMarkup {
	buttonsMap := k.initKeyboard()
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[LinksForUtil].Text, buttonsMap[LinksForUtil].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[ViewQueue].Text, buttonsMap[ViewQueue].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[ActualState].Text, buttonsMap[ActualState].ID),
		),
	)
}

func (k *Bot) keyboardAdmins() tgbotapi.InlineKeyboardMarkup {
	buttonsMap := k.initKeyboard()
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[CleanHistory].Text, buttonsMap[CleanHistory].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[RemoveFromQueue].Text, buttonsMap[RemoveFromQueue].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[LinksForUtil].Text, buttonsMap[LinksForUtil].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[ViewQueue].Text, buttonsMap[ViewQueue].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[ActualState].Text, buttonsMap[ActualState].ID),
		),
	)
}

type Quality int

const (
	FHD Quality = iota + 1
	_2K
	_4K
)

var qualities map[Quality]string = map[Quality]string{
	FHD: "1080",
	_2K: "1440",
	_4K: "2160",
}

var qualitiesInt map[Quality]int = map[Quality]int{
	FHD: 1080,
	_2K: 1440,
	_4K: 2160,
}

func (b *Bot) initKeyboardLink(link string) map[string]buttons {
	buttonsMap := make(map[string]buttons)

	buttonsMap[qualities[FHD]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[FHD], link),
		Text: qualities[FHD],
	}
	buttonsMap[qualities[_2K]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[_2K], link),
		Text: qualities[_2K],
	}
	buttonsMap[qualities[_4K]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[_4K], link),
		Text: qualities[_4K],
	}

	return buttonsMap
}

func (k *Bot) qualityKeyboard(link string) tgbotapi.InlineKeyboardMarkup {
	buttonsMap := k.initKeyboardLink(link)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[FHD]].Text, buttonsMap[qualities[FHD]].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[_2K]].Text, buttonsMap[qualities[_2K]].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[_4K]].Text, buttonsMap[qualities[_4K]].ID),
		),
	)
}
