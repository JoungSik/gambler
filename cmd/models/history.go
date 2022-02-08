package models

import "time"

type History struct {
	ID        uint `gorm:"primaryKey;auto_increment"`
	UserID    string
	Invest    int64
	Principal int64
	Result    int64
	Tax       int64
	Total     int64
	Diameter  int64
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
