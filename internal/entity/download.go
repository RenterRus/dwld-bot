package entity

import "time"

type TaskModel struct {
	Link      string
	UserID    string
	MessageID string
	ErrorMsg  string
	Quality   int
	SendAt    time.Time
}

type ServerModel struct {
	Name             string
	AllowedRootLinks []string
	Host             string
	Port             int
}

type LoadBy int

const (
	ByTime LoadBy = iota + 1
	ByUserID
	ByLink
	ByAny
)
