package dwld

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/RenterRus/dwld-bot/internal/entity"
	dwl "github.com/RenterRus/dwld-downloader/docs/proto/v1"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Идем по скачивальщикам

type DownloaderRepo struct {
	client     dwl.DownloaderClient
	cc         *grpc.ClientConn
	serverName string
}

func NewDWLD(host string, port int) DWLDModel {
	cc, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
	}

	return &DownloaderRepo{
		client: dwl.NewDownloaderClient(cc),
	}
}

func (d *DownloaderRepo) CleanHistory(ctx context.Context) ([]*entity.TaskRaw, error) {
	tasks, err := d.client.CleanHistory(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("CleanHistory: %w", err)
	}

	return lo.Map(tasks.History, taskToTaskRaw), nil
}

func (d *DownloaderRepo) DeleteFromQueue(ctx context.Context, link string) ([]*entity.TaskRaw, error) {
	tasks, err := d.client.DeleteFromQueue(ctx, &dwl.DeleteFromQueueRequest{
		Link: link,
	})
	if err != nil {
		return nil, fmt.Errorf("DeleteFromQueue: %w", err)
	}

	return lo.Map(tasks.LinksInWork, taskToTaskRaw), nil
}

func (d *DownloaderRepo) Queue(ctx context.Context) ([]*entity.TaskRaw, error) {
	tasks, err := d.client.Queue(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("Queue: %w", err)
	}

	return lo.Map(tasks.Queue, taskToTaskRaw), nil
}

func (d *DownloaderRepo) SetToQueue(ctx context.Context, link string, targetQuality int32) ([]*entity.TaskRaw, error) {
	tasks, err := d.client.SetToQueue(ctx, &dwl.SetToQueueRequest{
		Link:       link,
		MaxQuality: pointer.To(targetQuality),
	})
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return lo.Map(tasks.LinksInWork, taskToTaskRaw), nil
}

func (d *DownloaderRepo) Status(ctx context.Context) (*entity.Status, error) {
	tasks, err := d.client.Status(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("Status: %w", err)
	}

	return &entity.Status{
		Tasks:   lo.Map(tasks.LinksInWork, onWorkToTaskInfo),
		Sensors: tasks.Sensors,
	}, nil
}

func (d *DownloaderRepo) SetName(name string) DWLDModel {
	d.serverName = name
	return d
}
func (d *DownloaderRepo) Name() string {
	return d.serverName
}
