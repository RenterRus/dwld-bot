package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/RenterRus/dwld-bot/internal/entity"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TIMEOUT_REFRESH_MSG = 17
	TIMEOUT_DELETE_MSG  = 5
)

type BotRepo struct {
	db     persistent.SQLRepo
	bot    *tgbotapi.BotAPI
	notify chan struct{}
}

func NewBotRepo(bot *tgbotapi.BotAPI, db persistent.SQLRepo) BotModel {
	return &BotRepo{
		bot:    bot,
		db:     db,
		notify: make(chan struct{}, 1),
	}
}

func (r *BotRepo) DeleteMsg(userID, messageID string) error {
	user := 0
	message := 0
	var err error

	if user, err = strconv.Atoi(userID); err != nil {
		fmt.Printf("DeleteMsg.ParseInt(userID): %s\n", err.Error())
		return fmt.Errorf("DeleteMsg.ParseInt(userID): %w", err)
	}

	if message, err = strconv.Atoi(messageID); err != nil {
		fmt.Printf("DeleteMsg.ParseInt(messageID): %s\n", err.Error())
		return fmt.Errorf("DeleteMsg.ParseInt(messageID): %w", err)
	}

	_, err = r.bot.Request(tgbotapi.NewDeleteMessage(int64(user), message))
	if err != nil {
		return fmt.Errorf("DeleteMsg: %w", err)
	}

	return nil
}

func (r *BotRepo) SendMessage(chatID, message string) {
	chat := 0
	var err error

	if chat, err = strconv.Atoi(chatID); err != nil {
		fmt.Printf("SendMessage.ParseInt(chatID): %s\n", err.Error())
		return
	}

	res, err := r.bot.Send(tgbotapi.NewMessage(int64(chat), message))
	if err != nil {
		fmt.Println("SendMessage (Send):", err.Error())
		return
	}

	r.SetToQueue(&TaskToDelete{
		ChatID:    res.Chat.ID,
		MessageID: res.MessageID,
		Deadline:  time.Now().Add(time.Minute * TIMEOUT_DELETE_MSG),
	})

}

func (r *BotRepo) Processor() {
	t := time.NewTicker(time.Second * TIMEOUT_REFRESH_MSG)

	for {
		select {
		case <-r.notify:
			return
		case <-t.C:
			tasks := r.db.GetToDelete()
			for _, task := range tasks {
				fmt.Println("message to delete:", task.ChatID, task.MesssageID)

				err := r.DeleteMsg(task.ChatID, task.MesssageID)
				if err != nil {
					fmt.Println("DELETE MSG FAILED:", err.Error())
				}
				r.db.RemoveToDelete(&entity.ToDeleteTask{
					ChatID:     task.ChatID,
					MesssageID: task.MesssageID,
					DeleteAt:   task.DeleteAt,
				})

			}
		}
	}

}

func (r *BotRepo) Stop() {
	r.notify <- struct{}{}
}

func (r *BotRepo) SetToQueue(task *TaskToDelete) {
	r.db.SetToDelete(&entity.ToDeleteTask{
		ChatID:     strconv.Itoa(int(task.ChatID)),
		MesssageID: strconv.Itoa(task.MessageID),
		DeleteAt:   task.Deadline,
	})
}
