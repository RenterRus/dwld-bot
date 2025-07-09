package dwld

import (
	"context"

	"github.com/RenterRus/dwld-bot/internal/entity"
)

type DWLDModel interface {
	CleanHistory(ctx context.Context) ([]*entity.TaskRaw, error)
	DeleteFromQueue(ctx context.Context, link string) ([]*entity.TaskRaw, error)
	Queue(ctx context.Context) ([]*entity.TaskRaw, error)
	SetToQueue(ctx context.Context, link string, targetQuality int32) ([]*entity.TaskRaw, error)
	Status(ctx context.Context) (*entity.Status, error)
}
