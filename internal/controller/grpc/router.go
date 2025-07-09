package grpc

import (
	v1 "github.com/RenterRus/dwld-bot/internal/controller/grpc/apiv1"
	"github.com/RenterRus/dwld-bot/internal/usecase/bot"

	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(app *pbgrpc.Server, usecases bot.Bot) {
	v1.NewDownloadRoutes(app, usecases)
	reflection.Register(app)
}
