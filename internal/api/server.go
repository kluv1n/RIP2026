package api

import (
	"fmt"
	"net/http"
	"strings"

	"RIP2026/lab_materials/docs"
	"RIP2026/internal/app/config"
	"RIP2026/internal/app/handler"
	"RIP2026/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// StartServer godoc
// @title Battery Life API
// @version 1.0
// @description API для управления заявками расчета времени работы аккумуляторов
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func StartServer(cfg *config.Config, repo *repository.Repository) {
	h := handler.NewHandler(repo)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.ServicePort)
	docs.SwaggerInfo.BasePath = "/api"

	h.RegisterAPI(r)

	r.GET("/", h.GetBatteryTypes)
	r.GET("/battery/:id", h.GetBattery)
	r.GET("/battery-life/:id", h.GetBatteryLife)
	r.POST("/battery-life/add", h.AddBatteryToBatteryLife)
	r.POST("/battery-life/delete", h.DeleteBatteryLife)

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") || strings.HasPrefix(c.Request.URL.Path, "/users/") {
			c.JSON(http.StatusNotFound, gin.H{
				"description": "Not found",
			})
			return
		}
		c.Redirect(http.StatusFound, "/")
	})

	addr := fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort)
	logrus.Info("Server starting at ", addr)
	if err := r.Run(addr); err != nil {
		logrus.Fatal(err)
	}
}
