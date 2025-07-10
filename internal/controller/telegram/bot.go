package telegram

import (
	"context"

	"github.com/RenterRus/dwld-bot/internal/usecase/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// !!!
type BotConfig struct {
	Token         string
	AllowedChatID []string
	AdminChatID   []string
	BotUsecases   bot.Bot
}

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(conf BotConfig) BotModel {
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		panic(err)
	}

	return &Bot{
		bot: bot,
	}
}

func (b *Bot) Processor(ctx context.Context) {

}
