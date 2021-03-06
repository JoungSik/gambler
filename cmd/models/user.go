package models

import "time"

type User struct {
	ID        int64 `gorm:"primaryKey"`
	Email     string
	Name      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
