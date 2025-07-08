package grpc

import (
	v1 "dwld-bot/internal/controller/grpc/apiv1"
	"dwld-bot/internal/usecase/download"

	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(app *pbgrpc.Server, usecases download.Downloader) {
	v1.NewDownloadRoutes(app, usecases)
	reflection.Register(app)
}
