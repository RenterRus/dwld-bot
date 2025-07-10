package bot

import (
	"fmt"

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

func (b *BotCases) SetTask(task entity.TaskModel) error {
	if err := b.db.SetTask(task); err != nil {
		return fmt.Errorf("SetTask: %w", err)
	}

	return nil
}

func (b *BotCases) ViewTasks() ([]*entity.TaskModel, error) {
	// !!!

	return nil, nil
}

func (b *BotCases) DeleteTask(link string) error {
	// !!!

	return nil
}

func (b *BotCases) Status() (*entity.TaskInfo, error) {
	// !!!
	return nil, nil
}

func (b *BotCases) StorageServer(server entity.ServerModel) error {
	if err := b.db.StorageServer(server); err != nil {
		return fmt.Errorf("StorageServer: %w", err)
	}

	return nil
}
