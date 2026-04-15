package main

import (
	"github.com/joho/godotenv"

	"RIP2026/internal/api"
	"RIP2026/internal/app/config"
	"RIP2026/internal/app/dsn"
	"RIP2026/internal/app/repository"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	repo, err := repository.New(dsn.FromEnv())
	if err != nil {
		panic(err)
	}
	api.StartServer(cfg, repo)
}
