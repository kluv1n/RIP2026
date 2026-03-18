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
	creatorID := uint(h.Repository.GetCreatorID())
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

	batteryLives, err := h.Repository.GetBatteryLives(creatorID)
	if err != nil {
		logrus.Error(err)
	}

	draftID, hasDraft := h.Repository.GetDraftID(creatorID)
	draftIDInt := 0
	if hasDraft {
		draftIDInt = int(draftID)
	}
	cartCount := h.Repository.GetCartCount(creatorID)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"batteries":    batteries,
		"query":        searchQuery,
		"batteryLives": batteryLives,
		"draftID":      draftIDInt,
		"hasDraft":     hasDraft,
		"cartCount":    cartCount,
	})
}

func (h *Handler) GetBattery(ctx *gin.Context) {
	creatorID := uint(h.Repository.GetCreatorID())
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	battery, err := h.Repository.GetBattery(id)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	batteryLifeItem, _ := h.Repository.GetBatteryLifeForBattery(id, creatorID)
	hasInBatteryLife := batteryLifeItem != nil

	ctx.HTML(http.StatusOK, "battery.html", gin.H{
		"battery":          battery,
		"batteryLifeItem":  batteryLifeItem,
		"hasInBatteryLife": hasInBatteryLife,
		"batteryID":        id,
	})
}

func (h *Handler) GetBatteryLife(ctx *gin.Context) {
	creatorID := uint(h.Repository.GetCreatorID())
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	life, err := h.Repository.GetBatteryLife(id, creatorID)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.HTML(http.StatusOK, "battery_life.html", gin.H{
		"batteryLife":    life,
		"runtimeSummary": fmt.Sprintf("%.2f ч", life.TotalRuntimeHours),
		"isDraft":       life.Status == "draft",
	})
}

func (h *Handler) AddBatteryToBatteryLife(ctx *gin.Context) {
	creatorID := uint(h.Repository.GetCreatorID())
	batteryID, _ := strconv.Atoi(ctx.PostForm("battery_id"))
	currentMa, _ := strconv.Atoi(ctx.PostForm("current_ma"))
	if currentMa <= 0 {
		currentMa = 100
	}
	quantity, _ := strconv.Atoi(ctx.PostForm("quantity"))
	if quantity <= 0 {
		quantity = 1
	}
	if err := h.Repository.AddBatteryToBatteryLife(creatorID, batteryID, currentMa, quantity); err != nil {
		logrus.Error(err)
		if ctx.PostForm("return_to") == "battery" {
			ctx.Redirect(http.StatusFound, "/battery/"+strconv.Itoa(batteryID)+"?error=add")
		} else {
			ctx.Redirect(http.StatusFound, "/?error=add")
		}
		return
	}
	if ctx.PostForm("return_to") == "battery" {
		ctx.Redirect(http.StatusFound, "/battery/"+strconv.Itoa(batteryID)+"?added=1")
	} else {
		ctx.Redirect(http.StatusFound, "/?added=1")
	}
}

func (h *Handler) DeleteBatteryLife(ctx *gin.Context) {
	creatorID := uint(h.Repository.GetCreatorID())
	idStr := ctx.PostForm("battery_life_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	if err := h.Repository.DeleteBatteryLife(id, creatorID); err != nil {
		logrus.Error(err)
	}
	ctx.Redirect(http.StatusFound, "/")
}

// RegisterAPI регистрирует маршруты REST API под префиксом /api
func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")
	batteryTypes := api.Group("/battery_types")
	{
		batteryTypes.GET("", h.APIGetBatteryTypes)
		batteryTypes.GET("/:id", h.APIGetBatteryType)
		batteryTypes.POST("", h.APICreateBatteryType)
	}
	batteryLives := api.Group("/battery_lives")
	{
		batteryLives.GET("/cart", h.APIGetBatteryLifeCart)
		batteryLives.GET("", h.APIGetBatteryLives)
		batteryLives.GET("/:id", h.APIGetBatteryLife)
		batteryLives.PUT("/:id", h.APIEditBatteryLife)
		batteryLives.PUT("/:id/form", h.APIFormBatteryLife)
		batteryLives.PUT("/:id/finish", h.APIFinishBatteryLife)
		batteryLives.DELETE("/:id", h.APIDeleteBatteryLife)
	}
	batteryLifeItems := api.Group("/battery_life_items")
	{
		batteryLifeItems.POST("/add/:battery_type_id", h.APIAddToBatteryLife)
		batteryLifeItems.DELETE("/:battery_type_id/:battery_life_id", h.APIDeleteFromBatteryLife)
		batteryLifeItems.PUT("/:battery_type_id/:battery_life_id", h.APIEditInBatteryLife)
	}
	users := api.Group("/users")
	{
		users.POST("/register", h.APICreateUser)
		users.POST("/login", h.APISignIn)
		users.POST("/logout", h.APISignOut)
	}
}
