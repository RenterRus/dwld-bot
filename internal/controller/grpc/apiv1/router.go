package v1

import (
	proto "dwld-bot/docs/proto/v1"
	"dwld-bot/internal/usecase/download"

	pbgrpc "google.golang.org/grpc"
)

func NewDownloadRoutes(app *pbgrpc.Server, usecases download.Downloader) {
	r := &V1{
		u: usecases,
	}

	proto.RegisterDownloaderServer(app, r)
}
