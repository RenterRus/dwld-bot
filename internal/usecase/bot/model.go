package bot

import (
	"context"

	"github.com/RenterRus/dwld-bot/internal/entity"
)

type Bot interface {
	SetTask(entity.TaskModel) error
	ViewTasks(ctx context.Context, userID string) ([]*entity.TaskRaw, error)
	DeleteTask(ctx context.Context, link string)
	Status(ctx context.Context) ([]*entity.Status, error)

	StorageServer(server entity.ServerModel) error

	CleanHistory(ctx context.Context)
}
