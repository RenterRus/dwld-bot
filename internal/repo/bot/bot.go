package bot

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotRepo struct {
}

func NewBotRepo() BotModel {
	return &BotRepo{}
}

func (r *BotRepo) DeleteMsg(userID, messageID string) error {
	user := int64(0)
	message := int64(0)
	var err error

	if user, err = strconv.ParseInt(userID, 10, 64); err == nil {
		return fmt.Errorf("DeleteMsg.ParseInt(userID): %w", err)
	}

	if message, err = strconv.ParseInt(messageID, 10, 64); err == nil {
		return fmt.Errorf("DeleteMsg.ParseInt(userID): %w", err)
	}

	tgbotapi.NewDeleteMessage(user, int(message))
	return nil
}
