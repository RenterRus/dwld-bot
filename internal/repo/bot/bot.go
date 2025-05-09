package bot

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TIMEOUT_DELETE_MSG = 17
)

type BotRepo struct {
	bot    *tgbotapi.BotAPI
	notify chan struct{}
	tasks  map[int]*TaskToDelete
}

func NewBotRepo(bot *tgbotapi.BotAPI) BotModel {
	return &BotRepo{
		bot:    bot,
		notify: make(chan struct{}, 1),
		tasks:  make(map[int]*TaskToDelete),
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

func (r *BotRepo) Processor() {
	t := time.NewTicker(time.Second * TIMEOUT_DELETE_MSG)

	for {
		select {
		case <-r.notify:
			return
		case <-t.C:
			for _, task := range r.tasks {
				if task.Deadline.Unix() <= time.Now().Unix() {
					r.DeleteMsg(strconv.Itoa(int(task.ChatID)), strconv.Itoa(task.MessageID))
					delete(r.tasks, task.MessageID)
				}
			}
		}
	}

}

func (r *BotRepo) Stop() {
	for _, task := range r.tasks {
		r.DeleteMsg(strconv.Itoa(int(task.ChatID)), strconv.Itoa(task.MessageID))
	}
	r.notify <- struct{}{}
}

func (r *BotRepo) SetToQueue(task *TaskToDelete) {
	r.tasks[task.MessageID] = task
}
