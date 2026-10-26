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
	base
	indrive
	development
	music
	background
	within
	pacifier
	clipers
	curious
)

type Button struct {
	EngName string
	RuName  string
}

var qualities map[Folder]Button = map[Folder]Button{
	income: Button{
		EngName: "0",
		RuName:  "Без категории",
	},
	background: Button{
		EngName: "1",
		RuName:  "Ha фон-к",
	},
	base: Button{
		EngName: "2",
		RuName:  "Это база",
	},
	development: Button{
		EngName: "3",
		RuName:  "Развитие",
	},
	clipers: Button{
		EngName: "4",
		RuName:  "Клипсы",
	},
	pacifier: Button{
		EngName: "5",
		RuName:  "Пустышка",
	},
	indrive: Button{
		EngName: "6",
		RuName:  "B дорогу",
	},
	within: Button{
		EngName: "7",
		RuName:  "Совместно",
	},
	curious: Button{
		EngName: "8",
		RuName:  "Любопытно",
	},
	music: Button{
		EngName: "9",
		RuName:  "Музыка",
	},
}

func (b *Bot) initKeyboardLink(link string) map[string]buttons {
	buttonsMap := make(map[string]buttons)

	for pos, _ := range qualities {
		buttonsMap[qualities[pos].EngName] = buttons{
			ID:   fmt.Sprintf("%s|%s", qualities[pos].EngName, link),
			Text: qualities[pos].RuName,
		}
	}

	return buttonsMap
}

/*
	[Это база | base]		[В дорогу | indrive]	   	[Развитие | development]
	[Музыка | music]		[На фон-к | background]		[Совместно | within]
	[Пустышка | pacifier]	[Клипсы | clipers]	   		[Любопытно | curious]
*/

func (k *Bot) qualityKeyboard(link string) tgbotapi.InlineKeyboardMarkup {
	buttonsMap := k.initKeyboardLink(link)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[base].EngName].Text, buttonsMap[qualities[base].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[indrive].EngName].Text, buttonsMap[qualities[indrive].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[development].EngName].Text, buttonsMap[qualities[development].EngName].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[music].EngName].Text, buttonsMap[qualities[music].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[background].EngName].Text, buttonsMap[qualities[background].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[within].EngName].Text, buttonsMap[qualities[within].EngName].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[pacifier].EngName].Text, buttonsMap[qualities[pacifier].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[clipers].EngName].Text, buttonsMap[qualities[clipers].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[curious].EngName].Text, buttonsMap[qualities[curious].EngName].ID),
		),
	)
}
