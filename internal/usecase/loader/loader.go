package loader

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"time"

	"github.com/RenterRus/dwld-bot/internal/entity"
	"github.com/RenterRus/dwld-bot/internal/repo/bot"
	"github.com/RenterRus/dwld-bot/internal/repo/dwld"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
)

const (
	SEND_TIMEOUT = 71
)

type LoaderCase struct {
	db     persistent.SQLRepo
	bot    bot.BotModel
	notify chan struct{}
}

func NewLoader(db persistent.SQLRepo, bot bot.BotModel) Loader {
	return &LoaderCase{
		db:     db,
		bot:    bot,
		notify: make(chan struct{}, 1),
	}
}

func (l *LoaderCase) Stop() {
	l.notify <- struct{}{}
}

func (l *LoaderCase) Processor(ctx context.Context) {
	t := time.NewTicker(time.Second * SEND_TIMEOUT)

	for {
		select {
		case <-t.C:
			tasks, err := l.db.LoadTasks(entity.ByTime, entity.TaskModel{})
			if err != nil {
				fmt.Println("Processor.LoadTasks", err)
				break
			}

			hosts, err := l.db.LoadServers()
			if err != nil {
				fmt.Println("Processor.LoadServers", err)
				break
			}

			servers := make(map[string]*entity.ServerModel)

			for _, v := range hosts {
				servers[v.Name] = v
			}

			//  Идем по полученным ссылкам
			for _, v := range tasks {
				host, err := url.Parse(v.Link)
				if err != nil {
					// Почему-то не смогли распарсить ссылку. Пробуем пометить ошибку и идем к следующей ссылке
					fmt.Println("Processor.url.Parse", err)
					v.ErrorMsg = err.Error()
					l.db.SetTask(*v)
					continue
				}

				//  Идем по конфигам серверов и ищем подходящий
				for _, conf := range servers {
					// Разрешенных нет - идем к следующему серверу
					if !slices.Contains(conf.AllowedRootLinks, host.Host) && !slices.Contains(conf.AllowedRootLinks, "*") &&
						!slices.Contains(conf.AllowedRootLinks, "any") && !slices.Contains(conf.AllowedRootLinks, "all") {
						fmt.Println("A suitable downloader was not found.")
						continue
					}

					// Этот сервер подходит

					//  Разрешенные ссылки сервера содержат ссылку задачи
					_, err := dwld.NewDWLD(conf.Host, conf.Port).SetToQueue(ctx, v.Link, v.UserName, int32(v.Quality))
					if err != nil {
						fmt.Printf("Loader(NewDWLD): %s", err.Error())
						v.ErrorMsg = err.Error()
						// Если не смогли вставить ссылку в скачивальщик, то откидываем на 3 итераций вперед перед повтором (но вначале пройдет по остальным подходящим серверам)
						v.SendAt = time.Now().Add(time.Second * (3 * SEND_TIMEOUT))
						l.db.SetTask(*v)
						continue // Идем в следующий сервер
					}

					// Если отправка удалась - удаляем сообщение у пользователя (которое отправил сам бот)
					if err := l.bot.DeleteMsg(v.UserID, v.MessageID); err != nil {
						fmt.Printf("Loader(DeleteMsg): %s\n", err.Error())
						v.ErrorMsg = err.Error()
						l.db.SetTask(*v)
					}

					// Удаляем таску из базы, если удалось отправить в скачивальщик
					if err := l.db.DeleteTask(v.Link); err != nil {
						fmt.Printf("Loader(DeleteTask): %s", err.Error())
						v.ErrorMsg = err.Error()
						l.db.SetTask(*v)
					}

					fmt.Printf("Ссылка [%s] отправлена в скачивальщик %s с целевым качеством %d\n", v.Link, conf.Name, v.Quality)
					l.bot.SendMessage(v.UserID, fmt.Sprintf("Ссылка [%s] отправлена в скачивальщик %s с целевым качеством %d", v.Link, conf.Name, v.Quality))
					break
				}
			}

		case <-l.notify:
			fmt.Println("stop upload")
			return
		}
	}
}
