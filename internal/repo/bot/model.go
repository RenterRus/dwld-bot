package bot

import "time"

type TaskToDelete struct {
	ChatID    int64
	MessageID int
	Deadline  time.Time
}

type BotModel interface {
	DeleteMsg(userID, messageID string) error
	Processor()
	Stop()
	SetToQueue(*TaskToDelete)
	SendMessage(chatID, message string)
}
