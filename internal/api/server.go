package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"RIP2026/internal/app/config"
	"RIP2026/internal/app/handler"
	"RIP2026/internal/app/repository"
)

func StartServer(cfg *config.Config, repo *repository.Repository) {
	h := handler.NewHandler(repo)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", h.GetBatteryTypes)
	r.GET("/battery/:id", h.GetBattery)
	r.GET("/battery-life/:id", h.GetBatteryLife)
	r.POST("/battery-life/add", h.AddBatteryToBatteryLife)
	r.POST("/battery-life/delete", h.DeleteBatteryLife)

	addr := fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort)
	logrus.Info("Server starting at ", addr)
	if err := r.Run(addr); err != nil {
		logrus.Fatal(err)
	}
}
