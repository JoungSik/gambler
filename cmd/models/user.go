package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Server    string
	Name      string
	Amount    int64
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
