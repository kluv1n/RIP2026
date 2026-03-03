package main

import (
	"log"
	"RIP2026/internal/api"
)

func main() {
	log.Println("Battery life service start!")
	api.StartServer()
	log.Println("Battery life service terminated!")
}
