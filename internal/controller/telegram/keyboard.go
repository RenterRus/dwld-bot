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
	def Folder = iota
	one
	two
	three
	four
	five
	six
	seven
	eith
	nine
)

type Button struct {
	EngName string
	RuName  string
}

/*
[Полезно | one]			[На работу | two]	   	[Топ | three]
[В дорогу | four]		[Совместно | five]		[На фон-к | six]
[Trash | seven]			[На тренировку | eith]	 [Интересно | nine]
*/
var qualities map[Folder]Button = map[Folder]Button{
	def: Button{
		EngName: "0",
		RuName:  "Без категории",
	},
	one: Button{
		EngName: "1",
		RuName:  "Полезно",
	},
	two: Button{
		EngName: "2",
		RuName:  "На работу",
	},
	three: Button{
		EngName: "3",
		RuName:  "Топ",
	},
	four: Button{
		EngName: "4",
		RuName:  "В дорогу",
	},
	five: Button{
		EngName: "5",
		RuName:  "Совместно",
	},
	six: Button{
		EngName: "6",
		RuName:  "На фон-к",
	},
	seven: Button{
		EngName: "7",
		RuName:  "Trash",
	},
	eith: Button{
		EngName: "8",
		RuName:  "На тренировку",
	},
	nine: Button{
		EngName: "9",
		RuName:  "Интересно",
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

func (k *Bot) qualityKeyboard(link string) tgbotapi.InlineKeyboardMarkup {
	buttonsMap := k.initKeyboardLink(link)
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[one].EngName].Text, buttonsMap[qualities[one].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[two].EngName].Text, buttonsMap[qualities[two].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[three].EngName].Text, buttonsMap[qualities[three].EngName].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[four].EngName].Text, buttonsMap[qualities[four].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[five].EngName].Text, buttonsMap[qualities[five].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[six].EngName].Text, buttonsMap[qualities[six].EngName].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[seven].EngName].Text, buttonsMap[qualities[seven].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[eith].EngName].Text, buttonsMap[qualities[eith].EngName].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap[qualities[nine].EngName].Text, buttonsMap[qualities[nine].EngName].ID),
		),
	)
}
