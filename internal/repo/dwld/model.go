package dwld

import (
	"context"
)

type Task struct {
	Link          string
	Status        string
	TargetQuality string
	Name          *string
	Message       *string
}

type TaskInfo struct {
	Link           string
	Filename       string
	MoveTo         string
	TargetQuantity int64
	Procentage     float64
	Status         string
	TotalSize      float64
	CurrentSize    float64
	Message        string
}

type Status struct {
	Sensors string
	Tasks   []*TaskInfo
}

type DWLDModel interface {
	CleanHistory(ctx context.Context) ([]*Task, error)
	DeleteFromQueue(ctx context.Context, link string) ([]*Task, error)
	Queue(ctx context.Context) ([]*Task, error)
	SetToQueue(ctx context.Context, link string, targetQuality int32) ([]*Task, error)
	Status(ctx context.Context) (*Status, error)
}
