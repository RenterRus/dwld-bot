package v1

import (
	proto "github.com/RenterRus/dwld-bot/docs/proto/v1"
	"github.com/RenterRus/dwld-bot/internal/usecase/bot"
)

type V1 struct {
	proto.BotServer

	u bot.Bot
}
