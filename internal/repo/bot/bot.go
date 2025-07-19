package bot

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TIMEOUT_REFRESH_MSG = 17
	TIMEOUT_DELETE_MSG  = 585
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
		ChatID:    int64(chat),
		MessageID: res.MessageID,
		Deadline:  time.Now().Add(TIMEOUT_DELETE_MSG * time.Second),
	})
}

func (r *BotRepo) Processor() {
	t := time.NewTicker(time.Second * TIMEOUT_REFRESH_MSG)

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
