package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/RenterRus/dwld-bot/internal/entity"
	rbot "github.com/RenterRus/dwld-bot/internal/repo/bot"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
	botusecase "github.com/RenterRus/dwld-bot/internal/usecase/bot"
	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DEFAULT_QUALITY     = 4320 // 8K
	DEFAULT_TIMEOUT     = 3
	FAST_DELETE_TIMEOUT = 1
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
	deleteMessage  rbot.BotModel
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
		deleteMessage:  rbot.NewBotRepo(bot, db),
	}
}

func (b *Bot) Bot() *tgbotapi.BotAPI {
	return b.bot
}

func (b *Bot) toAdmins(msg string) {
	admins := make([]int64, 0, len(b.adminChatID))
	for _, v := range b.adminChatID {
		if id, err := strconv.Atoi(v); err == nil {
			admins = append(admins, int64(id))
		}
	}

	for _, v := range admins {
		b.sendMessage(tgbotapi.NewMessage(v, msg))
	}
}

func (b *Bot) Processor(ctx context.Context) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	allowedChats := make(map[string]*UserState)
	for _, v := range b.allowedChatID {
		allowedChats[v] = &UserState{
			Mode:    DownloadMode,
			IsAdmin: false,
		}
	}

	for _, v := range b.adminChatID {
		allowedChats[v] = &UserState{
			Mode:    DownloadMode,
			IsAdmin: true,
		}
	}

	go func() {
		b.deleteMessage.Processor()
	}()
	defer func() {
		b.deleteMessage.Stop()
	}()

	go func() {
		b.toAdmins("Бот запущен")
	}()

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			// Убеждаемся, что пользователь из разрешенного пула
			var msg tgbotapi.MessageConfig
			chatID := fmt.Sprintf("%d", update.Message.Chat.ID)
			isLinkInsert := false

			b.deleteMessage.SetToQueue(&rbot.TaskToDelete{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
				Deadline:  time.Now().Add(time.Minute * FAST_DELETE_TIMEOUT),
			})

			if _, ok := allowedChats[chatID]; ok {
				// Этот блок должен идти до валидации на url, т.к. в очереди, теоретически, может оказаться вообще не ссылка (ручной ввод)
				// Если режим удаления
				if allowedChats[chatID].Mode == RemoveMode {
					allowedChats[chatID].Mode = DownloadMode

					b.botCase.DeleteTask(ctx, update.Message.Text)

					b.sendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ссылка [%s] удалена", update.Message.Text)))

					continue
				}

				// Не получилось обнаружить ссылку
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
					Link:      update.Message.Text,
					Quality:   DEFAULT_QUALITY,
					UserName:  update.Message.From.UserName,
					UserID:    strconv.Itoa(int(update.Message.Chat.ID)),
					MessageID: strconv.Itoa(update.Message.MessageID),
					ErrorMsg:  "",
					SendAt:    time.Now().Add(time.Minute * DEFAULT_TIMEOUT),
				}); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Не удалось вставить в очередь ссылку %s. Причина: %v", update.Message.Text, err.Error()))
				} else {
					isLinkInsert = true
					//Ссылка встала в очередь
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ссылка [%s] взята в работу", update.Message.Text))
					//Прикрепляем клавитуру выбора качества для конкретной ссылки
					msg.ReplyMarkup = b.qualityKeyboard(update.Message.Text)

					// Удаляем сообщение пользователя
					b.deleteMessage.SetToQueue(&rbot.TaskToDelete{
						ChatID:    update.Message.Chat.ID,
						MessageID: update.Message.MessageID,
						Deadline:  time.Now().Add(time.Second * DEFAULT_TIMEOUT),
					})
				}

			} else { // Если нет, то даем ответ о запрещенном доступе
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Доступ запрещен: %d", int(update.Message.Chat.ID)))
			}

			// Отправляем сообщение
			var mInfo tgbotapi.Message
			var err error
			if mInfo, err = b.bot.Send(msg); err != nil {
				fmt.Println("Send", err)
			}
			if isLinkInsert {
				b.botCase.SetTask(entity.TaskModel{
					Link:      update.Message.Text,
					Quality:   DEFAULT_QUALITY,
					UserName:  update.Message.From.UserName,
					UserID:    strconv.Itoa(int(update.Message.Chat.ID)),
					MessageID: strconv.Itoa(mInfo.MessageID),
					ErrorMsg:  "",
					SendAt:    time.Now().Add(time.Minute * DEFAULT_TIMEOUT),
				})
			} else {
				b.deleteMessage.SetToQueue(&rbot.TaskToDelete{
					ChatID:    update.Message.Chat.ID,
					MessageID: mInfo.MessageID,
					Deadline:  time.Now().Add(time.Minute * DEFAULT_TIMEOUT),
				})
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
					sensors := strings.Builder{}
					sensors.WriteString("SENSORS\n")

					queues := strings.Builder{}
					queues.WriteString("QUEUE (CACHE)\n")
					b.sendMessage(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Тут отображаются очереди и показания датчиков скачивальщиков. Кнопка \"Показать очередь\" покажет очередь по всей системе с включением очереди бота"))

					for i, v := range status {
						sensors.WriteString(fmt.Sprintf("%d of %d: %s\n", (i + 1), len(status), v.ServerName))
						sensors.WriteString(v.Sensors)
						sensors.WriteString("\n")

						for n, task := range v.Tasks {
							if n > 0 {
								queues.WriteString("\n")
							}
							queues.WriteString(fmt.Sprintf("<b>G%d|L%d</b>: [%s][%d][%.2f][%.2f/%.2f][ %s ][%s] TO: [%s] %s\n", (i + 1), (n + 1), task.Status, int(task.TargetQuantity),
								task.Procentage, task.CurrentSize, task.TotalSize, task.Link, task.Filename, task.MoveTo, task.Message))
						}

						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("%s\n%s", sensors.String(), queues.String()))
						msg.ParseMode = tgbotapi.ModeHTML
						b.sendMessage(msg)

						sensors.Reset()
						queues.Reset()
					}

				}
			case ViewQueue:
				queue, err := b.botCase.ViewTasks(ctx, strconv.Itoa(int(update.CallbackQuery.Message.Chat.ID)))
				if err != nil {
					msg.Text = fmt.Sprintf("Ошибка получения всей очереди: %s", err.Error())
				} else {
					if len(queue) == 0 {
						b.sendMessage(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Все очереди свободны"))
					}

					resp := strings.Builder{}
					lastStatus := ""
					for _, v := range queue {
						if lastStatus != v.Status {
							resp.WriteString("\n")
						}

						lastStatus = v.Status
						resp.WriteString(fmt.Sprintf("[<b>%s</b>][%s][%s][%s][%s]\n", v.Status, v.TargetQuality, pointer.Get(v.Name), v.Link, pointer.Get(v.Message)))
					}
					msg.Text = resp.String()
					msg.ParseMode = tgbotapi.ModeHTML
				}

			default:
				data := strings.Split(update.CallbackData(), "|")
				if len(data) == 2 {
					task := entity.TaskModel{
						Link:      data[1],
						UserID:    strconv.Itoa(int(update.CallbackQuery.Message.Chat.ID)),
						MessageID: strconv.Itoa(update.CallbackQuery.Message.MessageID),
						ErrorMsg:  "",
						UserName:  update.CallbackQuery.From.UserName,
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

			b.sendMessage(msg)
		}

	}
}

func (b *Bot) sendMessage(c tgbotapi.Chattable) {
	var mInfo tgbotapi.Message
	var err error
	if mInfo, err = b.bot.Send(c); err != nil {
		fmt.Println("NewMessage", err)
	} else {
		b.deleteMessage.SetToQueue(&rbot.TaskToDelete{
			ChatID:    mInfo.Chat.ID,
			MessageID: mInfo.MessageID,
			Deadline:  time.Now().Add(time.Minute * DEFAULT_TIMEOUT),
		})
	}
}
