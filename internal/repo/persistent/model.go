package persistent

import (
	"github.com/RenterRus/dwld-bot/internal/entity"
)

type Task interface {
	SetTask(task entity.TaskModel) error
	LoadTasks(by entity.LoadBy, task entity.TaskModel) ([]*entity.TaskModel, error)
	DeleteTask(link string) error
}

type Server interface {
	StorageServer(server entity.ServerModel) error
	LoadServers() ([]*entity.ServerModel, error)
}

type SQLRepo interface {
	Task
	Server
}

type ServerDTO struct {
	Name             string `sql:"name"`
	AllowedRootLinks string `sql:"allowedRootLinks"`
	Host             string `sql:"host"`
	Port             int    `sql:"port"`
}

type TaskDTO struct {
	Link      string `sql:"link"`
	UserID    string `sql:"userID"`
	MessageID string `sql:"messageID"`
	ErrorMsg  string `sql:"errorMsgw"`
	Quality   int    `sql:"quality"`
	SendAt    string `sql:"sendingAt"`
	UserName  string `sql:"userName"`
}
