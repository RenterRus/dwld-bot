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
	db  persistent.SQLRepo
	bot bot.BotModel
}

func NewLoader(db persistent.SQLRepo, dwld dwld.DWLDModel, bot bot.BotModel) Loader {
	return &LoaderCase{
		db:  db,
		bot: bot,
	}
}

func (l *LoaderCase) Processor(ctx context.Context) {
	t := time.NewTicker(time.Second * SEND_TIMEOUT)

	for {
		select {
		case <-t.C:
			tasks, err := l.db.LoadTasks(entity.ByTime, entity.TaskModel{})
			if err != nil {
				fmt.Println("Processor.LoadTasks", err)
				continue
			}

			hosts, err := l.db.LoadServers()
			if err != nil {
				fmt.Println("Processor.LoadServers", err)
				continue
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
					if !slices.Contains(conf.AllowedRootLinks, host.Host) {
						continue
					}

					// Этот сервер подходит

					//  Разрешенные ссылки сервера содержат ссылку задачи
					_, err := dwld.NewDWLD(conf.Host, conf.Port).SetToQueue(ctx, v.Link, int32(v.Quality))
					if err != nil {
						v.ErrorMsg = err.Error()
						// Если не смогли вставить ссылку в скачивальщик, то откидываем на 3 итераций вперед перед повтором (но вначале пройдет по остальным подходящим серверам)
						v.SendAt = time.Now().Add(time.Second * (3 * SEND_TIMEOUT))
						l.db.SetTask(*v)
						continue // Идем в следующий сервер
					}

					// Если отправка удалась - удаляем сообщение у пользователя (которое отправил сам бот)
					if err := l.bot.DeleteMsg(v.UserID, v.MessageID); err != nil {
						v.ErrorMsg = err.Error()
						l.db.SetTask(*v)
					}

					// Удаляем таску из базы, если удалось отправить в скачивальщик
					if err := l.db.DeleteTask(v.Link); err != nil {
						v.ErrorMsg = err.Error()
						l.db.SetTask(*v)
					}

					// Раз мы оказались тут, то запись была отправлена на сервер, дальше отправлять ее не будем, идем к следующей ссылке
					break
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
