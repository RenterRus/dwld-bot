package v1

import (
	proto "github.com/RenterRus/dwld-bot/docs/proto/v1"
	"github.com/RenterRus/dwld-bot/internal/usecase/download"
)

type V1 struct {
	proto.DownloaderServer

	u download.Downloader
}
