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

type Folder int

const (
	income Folder = iota
	background
	trash
	learn
	music
	interesting
	within
)

var qualities map[Folder]string = map[Folder]string{
	income:      "inbox",
	background:  "Нафонк",
	trash:       "Нутакоэ",
	learn:       "Облучение",
	music:       "Клипсы",
	interesting: "Интересное",
	within:      "Совместно",
}

func (b *Bot) initKeyboardLink(link string) map[string]buttons {
	buttonsMap := make(map[string]buttons)

	buttonsMap[qualities[background]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[background], link),
		Text: qualities[background],
	}
	buttonsMap[qualities[trash]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[trash], link),
		Text: qualities[trash],
	}
	buttonsMap[qualities[learn]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[learn], link),
		Text: qualities[learn],
	}

	buttonsMap[qualities[music]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[music], link),
		Text: qualities[music],
	}
	buttonsMap[qualities[interesting]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[interesting], link),
		Text: qualities[interesting],
	}
	buttonsMap[qualities[within]] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[within], link),
		Text: qualities[within],
	}

	return buttonsMap
}

func (k *Bot) qualityKeyboard(link string) tgbotapi.InlineKeyboardMarkup {
	buttonsMap := k.initKeyboardLink(link)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[background]].Text, buttonsMap[qualities[background]].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[trash]].Text, buttonsMap[qualities[trash]].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[learn]].Text, buttonsMap[qualities[learn]].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[music]].Text, buttonsMap[qualities[music]].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[interesting]].Text, buttonsMap[qualities[interesting]].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[within]].Text, buttonsMap[qualities[within]].ID),
		),
	)
}
