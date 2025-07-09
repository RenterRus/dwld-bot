package persistent

import (
	"dwld-bot/pkg/sqldb"
)

// !!! Ходим по локальной бд

type persistentRepo struct {
	db *sqldb.DB
}

func NewSQLRepo(db *sqldb.DB) SQLRepo {
	return &persistentRepo{
		db: db,
	}
}

func (p *persistentRepo) Select(q string) ([]LinkModel, error) {
	return nil, nil
}

func (p *persistentRepo) Upsert(LinkModel) ([]LinkModel, error) {
	return nil, nil
}
