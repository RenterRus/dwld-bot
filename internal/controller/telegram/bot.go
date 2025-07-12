package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/RenterRus/dwld-bot/internal/entity"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
	botusecase "github.com/RenterRus/dwld-bot/internal/usecase/bot"
	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DEFAULT_QUALITY = 4320 // 8K
)

type BotConfig struct {
	Token         string
	AllowedChatID []string
	AdminChatID   []string
}

type Bot struct {
	bot            *tgbotapi.BotAPI
	adminChatID    []string
	allowedChatID  []string
	defaultQuality int
	botCase        botusecase.Bot
}

func NewBot(conf BotConfig, db persistent.SQLRepo) BotModel {
	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		panic(err)
	}

	return &Bot{
		bot:            bot,
		defaultQuality: DEFAULT_QUALITY,
		adminChatID:    conf.AdminChatID,
		allowedChatID:  conf.AllowedChatID,
		botCase:        botusecase.NewBotUsecases(db),
	}
}

func (b *Bot) toAdmins(msg string) {
	admins := make([]int64, 0, len(b.adminChatID))
	for _, v := range b.adminChatID {
		if id, err := strconv.Atoi(v); err == nil {
			admins = append(admins, int64(id))
		}
	}

	for _, v := range admins {
		b.bot.Send(tgbotapi.NewMessage(v, msg))
	}
}

func (b *Bot) Processor(ctx context.Context) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	allowedChats := make(map[string]*UserState)
	for _, v := range b.adminChatID {
		allowedChats[v] = &UserState{
			Mode:    DownloadMode,
			IsAdmin: true,
		}
	}

	for _, v := range b.allowedChatID {
		allowedChats[v] = &UserState{
			Mode:    DownloadMode,
			IsAdmin: false,
		}
	}

	go func() {
		b.toAdmins("Бот запущен")
	}()

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			// Убеждаемся, что пользователь из разрешенного пула
			var msg tgbotapi.MessageConfig
			chatID := fmt.Sprintf("%d", update.Message.Chat.ID)
			if _, ok := allowedChats[chatID]; ok {
				// Этот блок должен идти до валидации на url, т.к. в очереди, теоретически, может оказаться вообще не ссылка (ручной ввод)
				// Если режим удаления
				if allowedChats[chatID].Mode == RemoveMode {
					allowedChats[chatID].Mode = DownloadMode

					b.botCase.DeleteTask(ctx, update.Message.Text)
					b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ссылка [%s] удалена", update.Message.Text)))
					continue
				}

				// Не получилось обновружить ссылку
				if err := validate.Var(update.Message.Text, "url"); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMSG(update.Message.Chat.ID))

					stat, ok := allowedChats[fmt.Sprintf("%d", int(update.Message.Chat.ID))]
					if ok {
						if stat.IsAdmin {
							msg.ReplyMarkup = b.keyboardAdmins()
						} else {
							msg.ReplyMarkup = b.keyboardDefault()
						}
					}

					// Это ссылка, но вставка не удалась
				} else if err := b.botCase.SetTask(entity.TaskModel{
					Link:    update.Message.Text,
					Quality: DEFAULT_QUALITY,

					UserID:    strconv.Itoa(int(update.Message.Chat.ID)),
					MessageID: strconv.Itoa(update.Message.MessageID),
					ErrorMsg:  "",
					SendAt:    time.Now().Add(time.Minute * 5),
				}); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Не удалось вставить в очередь ссылку %s. Причина: %v", update.Message.Text, err.Error()))
				} else {
					//Ссылка встала в очередь
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ссылка [%s] взята в работу", update.Message.Text))
					//Прикрепляем клавитуру выбора качества для конкретной ссылки
					msg.ReplyMarkup = b.qualityKeyboard(update.Message.Text)
				}

			} else { // Если нет, то даем ответ о запрещенном доступе
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Доступ запрещен: %d", int(update.Message.Chat.ID)))
			}

			// Отправляем сообщение
			if _, err := b.bot.Send(msg); err != nil {
				fmt.Println("Send", err)
			}
		} else if update.CallbackQuery != nil { // Если пришел колбэк
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.bot.Request(callback); err != nil {
				fmt.Println("update.CallbackQuery", err)
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

			switch update.CallbackQuery.Data {
			case RemoveFromQueue:
				allowedChats[fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)].Mode = RemoveMode
				msg.Text = "Вставьте ссылку, которую надо удалить"

			case LinksForUtil:
				queue, err := b.botCase.ViewTasks(ctx, strconv.Itoa(int(update.CallbackQuery.Message.Chat.ID)))
				if err != nil {
					msg.Text = fmt.Sprintf("Ошибка получения всей очереди: %s", err.Error())
				} else {
					resp := strings.Builder{}
					for _, v := range queue {
						resp.WriteString(fmt.Sprintf("\"%s\",\n", v.Link))
					}
					msg.Text = resp.String()
				}

			case CleanHistory:
				b.botCase.CleanHistory(ctx)
				msg.Text = "Очистка истории отправлена в скачивальщики"

			case ActualState:
				status, err := b.botCase.Status(ctx)
				if err != nil {
					msg.Text = fmt.Sprintf("Ошибка получения актуального состояния: %s", err.Error())
				} else {
					b, err := json.Marshal(status)
					if err != nil {
						msg.Text = fmt.Sprintf("Ошибка обработки актуального состояния: %s", err.Error())
					}
					msg.Text = string(b)
				}
			case ViewQueue:
				queue, err := b.botCase.ViewTasks(ctx, strconv.Itoa(int(update.CallbackQuery.Message.Chat.ID)))
				if err != nil {
					msg.Text = fmt.Sprintf("Ошибка получения всей очереди: %s", err.Error())
				} else {
					resp := strings.Builder{}
					for _, v := range queue {
						resp.WriteString(fmt.Sprintf("[%s][%s][%s][%s][%s]\n", v.Status, v.TargetQuality, *v.Name, v.Link, *v.Message))
					}
					msg.Text = resp.String()
				}

			default:
				data := strings.Split(update.CallbackData(), "|")
				if len(data) == 2 {
					task := entity.TaskModel{
						Link:      data[1],
						UserID:    strconv.Itoa(int(update.CallbackQuery.Message.Chat.ID)),
						MessageID: strconv.Itoa(update.CallbackQuery.Message.MessageID),
						ErrorMsg:  "",
						Quality:   DEFAULT_QUALITY,
						SendAt:    time.Now().Add(time.Minute * 5),
					}

					switch data[0] {
					case qualities[FHD]:
						task.Quality = qualitiesInt[FHD]
					case qualities[_2K]:
						task.Quality = qualitiesInt[_2K]
					case qualities[_4K]:
						task.Quality = qualitiesInt[_4K]
					}

					if err := b.botCase.SetTask(task); err != nil {
						msg.Text = fmt.Sprintf("Не получилось вставить обновление, причина: %s", err.Error())
					}
					msg.Text = fmt.Sprintf("Для видео [%s] задано целевое качество [%d]", task.Link, task.Quality)
				} else {
					msg.Text = "Неожиданная команда"
				}
			}

			if _, err := b.bot.Send(msg); err != nil {
				fmt.Println("NewMessage", err)
			}
		}

	}
}
