package bot

type BotModel interface {
	DeleteMsg(userID, messageID string) error
}
