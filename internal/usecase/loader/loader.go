package loader

import (
	"github.com/RenterRus/dwld-bot/internal/repo/dwld"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
)

type LoaderCase struct {
	db   persistent.SQLRepo
	dwld dwld.DWLDModel
}

func NewLoader(db persistent.SQLRepo, dwld dwld.DWLDModel) Loader {
	return &LoaderCase{
		db:   db,
		dwld: dwld,
	}
}

func (l *LoaderCase) Processor() {
	// !!!
}
