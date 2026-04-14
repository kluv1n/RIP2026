package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"RIP2026/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

func (h *Handler) GetBatteryTypes(ctx *gin.Context) {
	var batteries []repository.BatteryType
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		batteries, err = h.Repository.GetBatteryTypes()
	} else {
		batteries, err = h.Repository.GetBatteryByTitle(searchQuery)
	}
	if err != nil {
		logrus.Error(err)
	}

	batteryLives, err := h.Repository.GetBatteryLives()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"batteries":    batteries,
		"query":        searchQuery,
		"batteryLives": batteryLives,
	})
}

func (h *Handler) GetBattery(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	battery, err := h.Repository.GetBattery(id)
	if err != nil {
		logrus.Error(err)
	}

	loadMa := battery.DemoLoadMa
	if loadMa <= 0 {
		loadMa = 100
	}
	hours := repository.RuntimeHours(battery.CapacityMah, loadMa)
	currentA := float64(loadMa) / 1000.0
	currentAStr := strings.Replace(fmt.Sprintf("%.2f", currentA), ".", ",", 1)
	runtimeHoursStr := strings.Replace(fmt.Sprintf("%.1f", hours), ".", ",", 1)

	ctx.HTML(http.StatusOK, "battery.html", gin.H{
		"battery":          battery,
		"batteryID":        id,
		"currentAStr":      currentAStr,
		"runtimeHoursStr":  runtimeHoursStr,
	})
}

func (h *Handler) GetBatteryLife(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	life, err := h.Repository.GetBatteryLife(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "battery_life.html", gin.H{
		"batteryLife":    life,
		"runtimeSummary": fmt.Sprintf("%.2f ч", life.TotalRuntimeHours),
	})
}
