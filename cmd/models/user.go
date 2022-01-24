package models

import "time"

const DEFAULT_AMOUNT = 100000

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	Amount    int64
	InitCount int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
