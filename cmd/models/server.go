package models

import "time"

type Server struct {
	ID        int64 `gorm:"primaryKey"`
	OwnerId   int64
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
