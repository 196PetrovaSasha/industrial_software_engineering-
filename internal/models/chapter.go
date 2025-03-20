package models

import (
	"time"
)

type Chapter struct {
	Id         int64 `gorm:"primary_key"`
	Name       string
	StartNode  int64
	Nodes      []int64             `gorm:"type:json;column:nodes"`
	Characters []int64             `gorm:"type:json;column:characters"`
	Status     int                 // 1 - черновик, 2 - на проверке, 3 - опубликована
	UpdatedAt  map[time.Time]int64 `gorm:"type:json;column:updated_at"`
	Author     int64
}
