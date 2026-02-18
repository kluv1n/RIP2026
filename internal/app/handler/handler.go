package handler

import (
	"fmt"
	"net/http"
	"strconv"

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

	applications, err := h.Repository.GetApplications()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"batteries":    batteries,
		"query":        searchQuery,
		"applications": applications,
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

	appItem, err := h.Repository.GetApplicationForBattery(id)
	hasInApplication := err == nil && appItem != nil

	ctx.HTML(http.StatusOK, "battery.html", gin.H{
		"battery":          battery,
		"appItem":          appItem,
		"hasInApplication": hasInApplication,
	})
}

func (h *Handler) GetApplication(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	app, err := h.Repository.GetApplication(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "application.html", gin.H{
		"app":            app,
		"runtimeSummary": fmt.Sprintf("%.2f ч", app.TotalRuntimeHours),
	})
}
