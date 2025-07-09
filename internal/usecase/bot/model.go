package bot

import "github.com/RenterRus/dwld-bot/internal/entity"

type Bot interface {
	SetTask(entity.TaskModel) error
	ViewTasks(userID string) ([]*entity.TaskModel, error)
	DeleteTask(link, userID string) error
	Status() (*entity.TaskInfo, error)

	StorageServer(server entity.ServerModel) error
}
