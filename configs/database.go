package configs

import (
	"github.com/JoungSik/gambler/cmd/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(debug bool) (*gorm.DB, error) {
	config := NewDBConfig(true)

	// dsn := "host=" + config.DB_HOST + " user=" + config.DB_USER + " password=" + config.DB_PASSWORD + " dbname=gambler port=5432 sslmode=disable"
	dsn := config.DB_USER + ":" + config.DB_PASSWORD + "@tcp(" + config.DB_HOST + ":3306)/gambler?charset=utf8mb4"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})

	return db, nil
}
