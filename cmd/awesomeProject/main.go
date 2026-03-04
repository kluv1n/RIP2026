package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"RIP2026/internal/app/config"
	"RIP2026/internal/app/dsn"
	"RIP2026/internal/app/repository"
	"RIP2026/internal/api"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("config: %v", err)
	}
	dsnStr := dsn.FromEnv()
	if dsnStr == "" {
		logrus.Fatal("DB_* env vars not set. Copy .env.example to .env and run docker-compose up -d postgres")
	}

	repo, err := repository.New(dsnStr)
	if err != nil {
		logrus.Fatalf("repository: %v", err)
	}
	logrus.Info("Battery life service start!")
	api.StartServer(cfg, repo)
}
