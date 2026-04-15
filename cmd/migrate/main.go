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
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
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
