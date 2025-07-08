package v1

import (
	proto "dwld-bot/docs/proto/v1"
	"dwld-bot/internal/usecase/download"
)

type V1 struct {
	proto.DownloaderServer

	u download.Downloader
}
