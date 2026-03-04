package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"RIP2026/internal/app/dsn"
	"RIP2026/internal/app/models"
)

func main() {
	_ = godotenv.Load()
	dsnStr := dsn.FromEnv()
	if dsnStr == "" {
		panic("set DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME in .env")
	}
	db, err := gorm.Open(postgres.Open(dsnStr), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	err = db.AutoMigrate(
		&models.User{},
		&models.BatteryType{},
		&models.BatteryLife{},
		&models.BatteryLifeItem{},
	)
	if err != nil {
		panic("migrate: " + err.Error())
	}
}
