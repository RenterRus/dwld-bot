package dwld

import (
	"github.com/RenterRus/dwld-bot/internal/entity"
	dwl "github.com/RenterRus/dwld-downloader/docs/proto/v1"
)

func taskToTaskRaw(t *dwl.Task, _ int) *entity.TaskRaw {
	return &entity.TaskRaw{
		Link:          t.Link,
		Status:        t.Status,
		TargetQuality: t.MaxQuantity,
		Name:          t.Name,
		Message:       t.Message,
	}
}

func onWorkToTaskInfo(t *dwl.OnWork, _ int) *entity.TaskInfo {
	return &entity.TaskInfo{
		Link:           t.Link,
		Filename:       t.Filename,
		MoveTo:         t.MoveTo,
		TargetQuantity: t.TargetQuantity,
		Procentage:     t.Procentage,
		Status:         t.Status,
		CurrentSize:    t.CurrentSize,
		TotalSize:      t.TotalSize,
		Message:        t.Message,
	}
}
