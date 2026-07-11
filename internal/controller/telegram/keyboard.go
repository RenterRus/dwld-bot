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
	indrive
	within
)

type Button struct {
	EngName string
	RuName  string
}

var qualities map[Folder]Button = map[Folder]Button{
	income: Button{
		EngName: "inbox",
		RuName:  "Без категории",
	},
	background: Button{
		EngName: "background",
		RuName:  "Ha фон-к",
	},
	trash: Button{
		EngName: "trash",
		RuName:  "Hy такоэ",
	},
	learn: Button{
		EngName: "learn",
		RuName:  "Облучение",
	},
	music: Button{
		EngName: "music",
		RuName:  "Клипсы",
	},
	interesting: Button{
		EngName: "interesting",
		RuName:  "ЛюбоПытное",
	},
	indrive: Button{
		EngName: "indrive",
		RuName:  "B дорогу",
	},
	within: Button{
		EngName: "within",
		RuName:  "Совместно",
	},
}

func (b *Bot) initKeyboardLink(link string) map[string]buttons {
	buttonsMap := make(map[string]buttons)

	buttonsMap[qualities[background].EngName] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[background].EngName, link),
		Text: qualities[background].RuName,
	}
	buttonsMap[qualities[trash].EngName] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[trash].EngName, link),
		Text: qualities[trash].RuName,
	}
	buttonsMap[qualities[learn].EngName] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[learn].EngName, link),
		Text: qualities[learn].RuName,
	}

	buttonsMap[qualities[music].EngName] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[music].EngName, link),
		Text: qualities[music].RuName,
	}
	buttonsMap[qualities[interesting].EngName] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[interesting].EngName, link),
		Text: qualities[interesting].RuName,
	}
	buttonsMap[qualities[within].EngName] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[within].EngName, link),
		Text: qualities[within].RuName,
	}

	buttonsMap[qualities[indrive].EngName] = buttons{
		ID:   fmt.Sprintf("%s|%s", qualities[within].EngName, link),
		Text: qualities[indrive].RuName,
	}

	return buttonsMap
}

func (k *Bot) qualityKeyboard(link string) tgbotapi.InlineKeyboardMarkup {
	buttonsMap := k.initKeyboardLink(link)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[income].EngName].Text, buttonsMap[qualities[income].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[trash].EngName].Text, buttonsMap[qualities[trash].EngName].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[learn].EngName].Text, buttonsMap[qualities[learn].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[music].EngName].Text, buttonsMap[qualities[music].EngName].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[indrive].EngName].Text, buttonsMap[qualities[indrive].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[background].EngName].Text, buttonsMap[qualities[background].EngName].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[interesting].EngName].Text, buttonsMap[qualities[interesting].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[within].EngName].Text, buttonsMap[qualities[within].EngName].ID),
		),
	)
}
