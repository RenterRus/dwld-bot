package bot

import (
	"context"

	"github.com/RenterRus/dwld-bot/internal/entity"
)

type Bot interface {
	SetTask(entity.TaskModel) error
	ViewTasks(ctx context.Context, userID string) ([]*entity.TaskRaw, error)
	DeleteTask(ctx context.Context, link string) error
	Status(ctx context.Context) (*entity.TaskInfo, error)

	StorageServer(server entity.ServerModel) error
}
