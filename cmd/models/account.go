package models

import "time"

const DEFAULT_AMOUNT = 100000

type Account struct {
	ID        uint `gorm:"primaryKey;auto_increment"`
	UserId    int64
	ServerId  int64
	Amount    int64
	InitCount int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
