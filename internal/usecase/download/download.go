package download

import (
	"dwld-bot/internal/repo/persistent"
	"dwld-bot/internal/repo/temporary"
	"dwld-bot/internal/usecase"
)

type downlaoder struct {
	dbRepo    *persistent.SQLRepo
	cacheRepo *temporary.Cache
}

func NewDownload(dbRepo *persistent.SQLRepo, cache *temporary.Cache) usecase.Downloader {
	return &downlaoder{
		dbRepo:    dbRepo,
		cacheRepo: cache,
	}
}
