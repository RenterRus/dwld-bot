package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/RenterRus/dwld-bot/internal/controller/grpc"
	"github.com/RenterRus/dwld-bot/internal/controller/telegram"
	rbot "github.com/RenterRus/dwld-bot/internal/repo/bot"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
	"github.com/RenterRus/dwld-bot/internal/usecase/bot"
	"github.com/RenterRus/dwld-bot/internal/usecase/loader"
	"github.com/RenterRus/dwld-bot/pkg/botserver"
	"github.com/RenterRus/dwld-bot/pkg/grpcserver"
	"github.com/RenterRus/dwld-bot/pkg/sqldb"
)

func NewApp(configPath string) error {
	lastSlash := 0
	for i, v := range configPath {
		if v == '/' {
			lastSlash = i
		}
	}

	conf, err := ReadConfig(configPath[:lastSlash], configPath[lastSlash+1:])
	if err != nil {
		return fmt.Errorf("ReadConfig: %w", err)
	}

	dbconn := sqldb.NewDB(conf.DB.PathToDB, conf.DB.NameDB)
	db := persistent.NewSQLRepo(dbconn)

	botsrv := botserver.NewBot(telegram.BotConfig{
		Token:         conf.TG.Token,
		AllowedChatID: conf.TG.AllowedIDs,
		AdminChatID:   conf.TG.Admins,
	}, db)

	upload := loader.NewLoader(db, rbot.NewBotRepo(botsrv.Bot(), db))

	downloadUsecases := bot.NewBotUsecases(db, upload)

	go func() {
		upload.Processor(context.Background())
	}()
	go func() {
		botsrv.Start()
	}()

	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(conf.GRPC.Host, strconv.Itoa(conf.GRPC.Port)))
	grpc.NewRouter(grpcServer.App, downloadUsecases)
	grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal: %s\n", s.String())
	case err = <-grpcServer.Notify():
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	botsrv.Stop()
	upload.Stop()
	dbconn.Close()
	err = grpcServer.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	return nil
}
