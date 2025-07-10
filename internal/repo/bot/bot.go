package bot

type BotRepo struct {
}

func NewBotRepo() BotModel {
	return &BotRepo{}
}

func (r *BotRepo) DeleteMsg(userID, messageID string) error {
	// !!!
	return nil
}
