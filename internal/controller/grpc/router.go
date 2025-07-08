package grpc

import (
	v1 "dwld-bot/internal/controller/grpc/apiv1"
	"dwld-bot/internal/usecase"

	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(app *pbgrpc.Server, usecases usecase.Downloader) {
	v1.NewDownloadRoutes(app, usecases)
	reflection.Register(app)
}
