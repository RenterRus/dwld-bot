package bot

import (
	"github.com/RenterRus/dwld-bot/internal/entity"
	"github.com/RenterRus/dwld-bot/internal/repo/dwld"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
)

type BotCases struct {
	db   persistent.SQLRepo
	dwld dwld.DWLDModel
}

func NewDownloadUsecases(db persistent.SQLRepo, dwld dwld.DWLDModel) Bot {
	return &BotCases{
		db:   db,
		dwld: dwld,
	}
}

func (b *BotCases) SetTask(entity.TaskModel) error {
	// !!!
	return nil
}

func (b *BotCases) ViewTasks(userID string) ([]*entity.TaskModel, error) {
	// !!!
	return nil, nil
}

func (b *BotCases) DeleteTask(link, userID string) error {
	// !!!
	return nil
}

func (b *BotCases) Status() (*entity.TaskInfo, error) {
	// !!!
	return nil, nil
}

func (b *BotCases) StorageServer(server entity.ServerModel) error {
	// !!!
	return nil
}
