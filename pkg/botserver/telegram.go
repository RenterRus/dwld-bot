package botserver

import (
	"context"
	"fmt"
	"time"

	"github.com/RenterRus/dwld-bot/internal/controller/telegram"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Поднимаем бота
const (
	LIFETIME = 1
)

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
	var ctx context.Context
	var cncl context.CancelFunc
	ticker := time.NewTicker(time.Hour * LIFETIME)

	go func() {
		for {
			fmt.Println("Bot starting")
			ctx, cncl = context.WithCancel(context.Background())
			go b.bot.Processor(ctx)
			<-ticker.C
			fmt.Println("Bot stopping")
			cncl()
			time.Sleep(time.Second * 5)
		}
	}()

	<-b.notify
	cncl()
	time.Sleep(time.Second * 3)
}

func (b *Bot) Stop() {
	b.notify <- struct{}{}
}
