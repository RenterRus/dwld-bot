package v1

import (
	proto "dwld-bot/docs/proto/v1"
	"dwld-bot/internal/usecase"
)

type V1 struct {
	proto.DownloaderServer

	u usecase.Downloader
}
