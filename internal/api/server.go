package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"RIP2026/internal/app/handler"
	"RIP2026/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", h.GetBatteryTypes)
	r.GET("/battery/:id", h.GetBattery)
	r.GET("/battery-life/:id", h.GetBatteryLife)

	r.Run()
	log.Println("Server down")
}
