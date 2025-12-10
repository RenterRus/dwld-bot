package persistent

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RenterRus/dwld-bot/internal/entity"
	"github.com/RenterRus/dwld-bot/pkg/sqldb"
)

// Идем в локальную БД

type persistentRepo struct {
	db *sqldb.DB
}

func NewSQLRepo(db *sqldb.DB) SQLRepo {
	return &persistentRepo{
		db: db,
	}
}

func (p *persistentRepo) SetTask(task entity.TaskModel) error {
	if _, err := p.db.Exec("insert into links(link, quality, sendingAt, userID, messageID, errorMsg, userName) values ($1, $2, $3, $4, $5, $6, $7) on conflict (link) do update set quality = excluded.quality, sendingAt = excluded.sendingAt, userID = excluded.userID, messageID = excluded.messageID, errorMsg = excluded.errorMsg, userName = excluded.userName;",
		task.Link, task.Quality, task.SendAt.Unix(), task.UserID, task.MessageID, task.ErrorMsg, task.UserName); err != nil {

		return fmt.Errorf("SetTask: %w", err)
	}

	return nil
}

func (p *persistentRepo) LoadTasks(by entity.LoadBy, task entity.TaskModel) ([]*entity.TaskModel, error) {
	var rows *sql.Rows
	var err error

	switch by {
	case entity.ByLink:
		rows, err = p.db.Select("select link, quality, sendingAt, userID, messageID, errorMsg, userName from links where link = $1", task.Link)
		if err != nil {
			return nil, fmt.Errorf("LoadTask(Select by links): %w", err)
		}
	case entity.ByAny:
		rows, err = p.db.Select("select link, quality, sendingAt, userID, messageID, errorMsg, userName from links")
		if err != nil {
			return nil, fmt.Errorf("LoadTask(Select by any): %w", err)
		}
	case entity.ByTime:
		rows, err = p.db.Select("select link, quality, sendingAt, userID, messageID, errorMsg, userName from links where sendingAt < $1", strconv.Itoa(int(time.Now().Unix())))
		if err != nil {
			return nil, fmt.Errorf("LoadTask(Select by time): %w", err)
		}
	case entity.ByUserID:
		fallthrough
	default:
		rows, err = p.db.Select("select link, quality, sendingAt, userID, messageID, errorMsg, userName from links where userID = $1", task.UserID)
		if err != nil {
			return nil, fmt.Errorf("LoadTask(Select by userID/default): %w", err)
		}
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("LoadTask: ", err.Error())
		}
	}()

	resp := make([]*entity.TaskModel, 0)
	var row TaskDTO
	for rows.Next() {
		err := rows.Scan(&row.Link, &row.Quality, &row.SendAt, &row.UserID, &row.MessageID, &row.ErrorMsg, &row.UserName)
		if err != nil {
			return nil, fmt.Errorf("LoadTask(Scan): %w", err)
		}

		sendAt, err := strconv.Atoi(row.UserID)
		if err != nil {
			return nil, fmt.Errorf("LoadTasks(atoi): %w", err)
		}

		resp = append(resp, &entity.TaskModel{
			Link:      row.Link,
			UserID:    row.UserID,
			MessageID: row.MessageID,
			ErrorMsg:  row.ErrorMsg,
			Quality:   row.Quality,
			UserName:  row.UserName,
			SendAt:    time.Unix(int64(sendAt), 0),
		})
	}

	return resp, nil
}

func (p *persistentRepo) DeleteTask(link string) error {
	if _, err := p.db.Exec("delete from links where link = $1", link); err != nil {
		return fmt.Errorf("DeleteTask: %w", err)
	}

	return nil
}

func (p *persistentRepo) StorageServer(server entity.ServerModel) error {
	allowedRoots := strings.Builder{}
	for i, v := range server.AllowedRootLinks {
		if i > 0 {
			allowedRoots.WriteString(",")
		}
		allowedRoots.WriteString(v)
	}

	if _, err := p.db.Exec("insert into downloaders(name, allowedRootLinks, host, port) values ($1, $2, $3, $4) on conflict (name) do update set allowedRootLinks = excluded.allowedRootLinks, host = excluded.host, port = excluded.port;",
		server.Name, allowedRoots.String(), server.Host, server.Port); err != nil {
		return fmt.Errorf("SetStorage: %w", err)
	}

	return nil
}

func (p *persistentRepo) LoadServers() ([]*entity.ServerModel, error) {
	rows, err := p.db.Select("select name, allowedRootLinks, host, port from downloaders")
	if err != nil {
		return nil, fmt.Errorf("SelectHistory: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("LoadTask: ", err.Error())
		}
	}()

	resp := make([]*entity.ServerModel, 0)
	var row ServerDTO
	for rows.Next() {
		err := rows.Scan(&row.Name, &row.AllowedRootLinks, &row.Host, &row.Port)
		if err != nil {
			return nil, fmt.Errorf("LoadServers(Scan): %w", err)
		}

		resp = append(resp, &entity.ServerModel{
			Name:             row.Name,
			AllowedRootLinks: strings.Split(row.AllowedRootLinks, ","),
			Host:             row.Host,
			Port:             row.Port,
		})
	}

	return resp, nil
}

func (p *persistentRepo) SetToDelete(t *entity.ToDeleteTask) {
	if _, err := p.db.Exec("insert into to_delete(chatID, messageID, deleteAt) values ($1, $2, $3)",
		t.ChatID, t.MesssageID, strconv.Itoa(int(t.DeleteAt.Unix()))); err != nil {
		fmt.Printf("SetToDelete: %s", err.Error())
	}

}
func (p *persistentRepo) GetToDelete() []*entity.ToDeleteTask {
	var rows *sql.Rows
	var err error

	rows, err = p.db.Select("select chatID, messageID, deleteAt from to_delete where deleteAt < $1", strconv.Itoa(int(time.Now().Unix())))
	if err != nil {
		fmt.Printf("GetToDelete (select): %s", err.Error())
		return nil
	}

	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println("LoadTask: ", err.Error())
		}
	}()

	resp := make([]*entity.ToDeleteTask, 0)
	var row ToDeleteTasDTO
	for rows.Next() {
		err := rows.Scan(&row.ChatID, &row.MesssageID, &row.DeleteAt)
		if err != nil {
			fmt.Printf("GetToDelete (scan): %s", err.Error())
			return nil
		}

		deleteAt, err := strconv.Atoi(row.DeleteAt)
		if err != nil {
			fmt.Printf("GetToDelete (scan): %s", err.Error())
			return nil
		}

		resp = append(resp, &entity.ToDeleteTask{
			ChatID:     row.ChatID,
			MesssageID: row.MesssageID,
			DeleteAt:   time.Unix(int64(deleteAt), 0),
		})
	}

	return resp
}

func (p *persistentRepo) RemoveToDelete(t *entity.ToDeleteTask) {
	if _, err := p.db.Exec("delete from to_delete where chatID = $1 and messageID = $2", t.ChatID, t.MesssageID); err != nil {
		fmt.Printf("RemoveToDelete: %s", err.Error())
	}
}
