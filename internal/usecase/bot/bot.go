package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/RenterRus/dwld-bot/internal/entity"
	"github.com/RenterRus/dwld-bot/internal/repo/dwld"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
)

type BotCases struct {
	db   persistent.SQLRepo
	dwld dwld.DWLDModel
}

func NewBotUsecases(db persistent.SQLRepo) Bot {
	return &BotCases{
		db: db,
	}
}

func (b *BotCases) servers() ([]dwld.DWLDModel, error) {
	servers, err := b.db.LoadServers()
	if err != nil {
		return nil, fmt.Errorf("servers: %w", err)
	}

	resp := make([]dwld.DWLDModel, 0, len(servers))

	for _, v := range servers {
		resp = append(resp, dwld.NewDWLD(v.Host, v.Port))
	}

	return resp, nil
}

func (b *BotCases) SetTask(task entity.TaskModel) error {
	if err := b.db.SetTask(task); err != nil {
		return fmt.Errorf("SetTask: %w", err)
	}

	return nil
}

func (b *BotCases) ViewTasks(ctx context.Context, userID string) ([]*entity.TaskRaw, error) {
	tasks, err := b.db.LoadTasks(entity.ByUserID, entity.TaskModel{
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("ViewTasks (LoadTasks): %w", err)
	}

	servers, err := b.servers()
	if err != nil {
		return nil, fmt.Errorf("ViewTasks (servers): %w", err)
	}

	resp := make([]*entity.TaskRaw, 0, len(tasks)*len(servers))
	for _, v := range tasks {
		resp = append(resp, &entity.TaskRaw{
			Link:          v.Link,
			Status:        "in queue",
			TargetQuality: strconv.Itoa(v.Quality),
			Name:          pointer.To("undefinded"),
			Message:       &v.ErrorMsg,
		})
	}

	for _, server := range servers {
		tasksRaw, err := server.Queue(ctx)
		if err != nil {
			fmt.Printf("ViewTasks (Queue): %s\n", err.Error())
		}
		resp = append(resp, tasksRaw...)
	}

	return resp, nil
}

func (b *BotCases) DeleteTask(ctx context.Context, link string) {
	if err := b.db.DeleteTask(link); err != nil {
		fmt.Printf("DeleteTask(DelteTask): %s\n", err.Error())
	}

	servers, err := b.servers()
	if err != nil {
		fmt.Printf("DeleteTask(servers): %s\n", err.Error())
	} else {
		for _, server := range servers {
			if _, err := server.DeleteFromQueue(ctx, link); err != nil {
				fmt.Printf("DeleteTask(DeleteFromQueue): %s\n", err.Error())
			}
		}
	}
}

func (b *BotCases) Status(ctx context.Context) ([]*entity.Status, error) {
	resp := make([]*entity.Status, 0, 2)

	servers, err := b.servers()
	if err != nil {
		return nil, fmt.Errorf("Status(servers): %w", err)
	} else {
		for _, server := range servers {
			var status *entity.Status
			if status, err = server.Status(ctx); err != nil {
				fmt.Printf("Status(Status): %s\n", err.Error())
			}
			resp = append(resp, &entity.Status{
				Sensors: pointer.Get(status).Sensors,
				Tasks:   pointer.Get(status).Tasks,
			})
		}
	}

	return resp, nil
}

func (b *BotCases) CleanHistory(ctx context.Context) {
	servers, err := b.servers()
	if err != nil {
		fmt.Printf("CleanHistory(servers): %s\n", err.Error())
	} else {
		for _, server := range servers {
			if _, err = server.CleanHistory(ctx); err != nil {
				fmt.Printf("CleanHistory(CleanHistory): %s\n", err.Error())
			}
		}
	}
}

func (b *BotCases) StorageServer(server entity.ServerModel) error {
	if err := b.db.StorageServer(server); err != nil {
		return fmt.Errorf("StorageServer: %w", err)
	}

	return nil
}
