package botserver

import (
	"context"

	"github.com/RenterRus/dwld-bot/internal/controller/telegram"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Поднимаем бота

type Bot struct {
	bot    telegram.BotModel
	notify chan struct{}
}

func NewBot(conf telegram.BotConfig, db persistent.SQLRepo) *Bot {
	return &Bot{
		bot:    telegram.NewBot(conf, db),
		notify: make(chan struct{}, 1),
	}

}

func (b *Bot) Bot() *tgbotapi.BotAPI {
	return b.bot.Bot()
}

func (b *Bot) Start() {
	ctx, cncl := context.WithCancel(context.Background())
	go func() {
		b.bot.Processor(ctx)
	}()

	<-b.notify
	cncl()
}

func (b *Bot) Stop() {
	b.notify <- struct{}{}
}
