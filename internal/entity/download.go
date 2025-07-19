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

type TaskInfo struct {
	Link           string
	Filename       string
	MoveTo         string
	TargetQuantity int64
	Procentage     float64
	Status         string
	TotalSize      float64
	CurrentSize    float64
	Message        string
}

type Status struct {
	ServerName string
	Sensors    string
	Tasks      []*TaskInfo
}

type TaskRaw struct {
	Link          string
	Status        string
	TargetQuality string
	Name          *string
	Message       *string
}
