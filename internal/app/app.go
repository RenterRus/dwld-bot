package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/RenterRus/dwld-bot/internal/controller/grpc"
	"github.com/RenterRus/dwld-bot/internal/usecase/download"
	"github.com/RenterRus/dwld-bot/pkg/grpcserver"
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

	// !!! db
	// !!! tg

	downloadUsecases := download.NewDownloadUsecases()

	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(strconv.Itoa(conf.GRPC.Port)))
	grpc.NewRouter(grpcServer.App, downloadUsecases)
	grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal: %s\n", s.String())
	case err = <-grpcServer.Notify():
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	return nil
}
