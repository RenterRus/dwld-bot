package persitent

import "dwld-bot/internal/entity"

type Task interface {
	SetTask(task entity.TaskModel) error
	LoadTasks(by entity.LoadBy, task entity.TaskModel) ([]*entity.TaskModel, error)
	DeleteTask(link string) error
}

type Server interface {
	StorageServer(server entity.ServerModel) error
	LoadServers(AllowedRootLinks []string) ([]*entity.ServerModel, error)
}

type SQLUsecase interface {
	Task
	Server
}
