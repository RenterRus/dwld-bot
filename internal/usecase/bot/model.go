package bot

import "github.com/RenterRus/dwld-bot/internal/entity"

type Bot interface {
	SetTask(entity.TaskModel) error
	ViewTasks() ([]*entity.TaskModel, error)
	DeleteTask(link string) error
	Status() (*entity.TaskInfo, error)

	StorageServer(server entity.ServerModel) error
}
