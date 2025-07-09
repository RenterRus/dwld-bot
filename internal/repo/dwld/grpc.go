package dwld

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	dwl "github.com/RenterRus/dwld-downloader/docs/proto/v1"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Идем по скачивальщикам

type Downloader struct {
	client dwl.DownloaderClient
}

func NewDWLD() DWLDModel {
	cc, err := grpc.NewClient("127.0.0.1:8999")
	if err != nil {
		fmt.Println(err)
	}
	return &Downloader{
		client: dwl.NewDownloaderClient(cc),
	}
}

func (d *Downloader) CleanHistory(ctx context.Context) ([]*Task, error) {
	tasks, err := d.client.CleanHistory(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("CleanHistory: %w", err)
	}

	return lo.Map(tasks.History, func(t *dwl.Task, _ int) *Task {
		return &Task{
			Link:          t.Link,
			Status:        t.Status,
			TargetQuality: t.MaxQuantity,
			Name:          t.Name,
			Message:       t.Message,
		}
	}), nil
}

func (d *Downloader) DeleteFromQueue(ctx context.Context, link string) ([]*Task, error) {
	tasks, err := d.client.DeleteFromQueue(ctx, &dwl.DeleteFromQueueRequest{
		Link: link,
	})
	if err != nil {
		return nil, fmt.Errorf("DeleteFromQueue: %w", err)
	}

	return lo.Map(tasks.LinksInWork, func(t *dwl.Task, _ int) *Task {
		return &Task{
			Link:          t.Link,
			Status:        t.Status,
			TargetQuality: t.MaxQuantity,
			Name:          t.Name,
			Message:       t.Message,
		}
	}), nil
}

func (d *Downloader) Queue(ctx context.Context) ([]*Task, error) {
	tasks, err := d.client.Queue(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("Queue: %w", err)
	}

	return lo.Map(tasks.Queue, func(t *dwl.Task, _ int) *Task {
		return &Task{
			Link:          t.Link,
			Status:        t.Status,
			TargetQuality: t.MaxQuantity,
			Name:          t.Name,
			Message:       t.Message,
		}
	}), nil
}

func (d *Downloader) SetToQueue(ctx context.Context, link string, targetQuality int32) ([]*Task, error) {
	tasks, err := d.client.SetToQueue(ctx, &dwl.SetToQueueRequest{
		Link:       link,
		MaxQuality: pointer.To(targetQuality),
	})
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return lo.Map(tasks.LinksInWork, func(t *dwl.Task, _ int) *Task {
		return &Task{
			Link:          t.Link,
			Status:        t.Status,
			TargetQuality: t.MaxQuantity,
			Name:          t.Name,
			Message:       t.Message,
		}
	}), nil
}

func (d *Downloader) Status(ctx context.Context) (*Status, error) {
	tasks, err := d.client.Status(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("Status: %w", err)
	}

	return &Status{
		Tasks: lo.Map(tasks.LinksInWork, func(t *dwl.OnWork, _ int) *TaskInfo {
			return &TaskInfo{
				Link:           t.Link,
				Filename:       t.Filename,
				MoveTo:         t.MoveTo,
				TargetQuantity: t.TargetQuantity,
				Procentage:     t.Procentage,
				Status:         t.Status,
				TotalSize:      t.TotalSize,
				CurrentSize:    t.CurrentSize,
				Message:        t.Message,
			}
		}),
		Sensors: tasks.Sensors,
	}, nil
}
