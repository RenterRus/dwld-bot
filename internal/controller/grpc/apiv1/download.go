package v1

import (
	"context"
	"fmt"

	v1 "github.com/RenterRus/dwld-bot/docs/proto/v1"
	"github.com/RenterRus/dwld-bot/internal/entity"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (v *V1) RegisterDownloader(ctx context.Context, in *v1.RegisterDownloaderRequest) (*emptypb.Empty, error) {
	if err := v.u.StorageServer(entity.ServerModel{
		Name:             in.ServerName,
		AllowedRootLinks: in.AllowedRootLinks,
		Host:             in.ServerHost,
		Port:             int(in.ServerPort),
	}); err != nil {
		return nil, fmt.Errorf("RegisterDownloaderv: %w", err)
	}

	return nil, nil
}
