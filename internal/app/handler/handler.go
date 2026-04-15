package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"RIP2026/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func batteryLifeStatusLabel(status string) string {
	switch status {
	case "draft", "черновик":
		return "Черновик"
	case "deleted", "удалён", "удален":
		return "Удалена"
	case "formed", "сформирован":
		return "Сформирована"
	case "completed", "завершён", "завершен":
		return "Завершена"
	case "rejected", "отклонён", "отклонен":
		return "Отклонена"
	default:
		if status == "" {
			return "—"
		}
		return status
	}
}

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

// demoLoadMaForCapacity — типовой ток (мА) для блока «Ток / время» на карточке; часы = мА·ч / мА.
func demoLoadMaForCapacity(capacityMah int) int {
	switch {
	case capacityMah >= 2800:
		return 500
	case capacityMah >= 2000:
		return 200
	case capacityMah >= 1400:
		return 250
	default:
		return 100
	}
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

	batteryLifeList, err := h.Repository.GetBatteryLives(creatorID)
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
		"batteryLifeList": batteryLifeList,
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

	draftID, hasDraft := h.Repository.GetDraftID(creatorID)
	draftIDInt := 0
	if hasDraft {
		draftIDInt = int(draftID)
	}
	cartCount := h.Repository.GetCartCount(creatorID)

	loadMa := demoLoadMaForCapacity(battery.CapacityMah)
	hours := repository.RuntimeHours(battery.CapacityMah, loadMa)
	currentA := float64(loadMa) / 1000.0
	currentAStr := strings.Replace(fmt.Sprintf("%.2f", currentA), ".", ",", 1)
	runtimeHoursStr := strings.Replace(fmt.Sprintf("%.1f", hours), ".", ",", 1)

	ctx.HTML(http.StatusOK, "battery.html", gin.H{
		"battery":           battery,
		"batteryLifeItem":   batteryLifeItem,
		"hasInBatteryLife":  hasInBatteryLife,
		"batteryID":         id,
		"currentAStr":       currentAStr,
		"runtimeHoursStr":   runtimeHoursStr,
		"query":             "",
		"draftID":           draftIDInt,
		"hasDraft":          hasDraft,
		"cartCount":         cartCount,
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

	isDraft := life.Status == "draft" || life.Status == "черновик"
	ctx.HTML(http.StatusOK, "battery_life.html", gin.H{
		"batteryLife":    life,
		"runtimeSummary": fmt.Sprintf("%.2f ч", life.TotalRuntimeHours),
		"isDraft":        isDraft,
		"statusLabel":    batteryLifeStatusLabel(life.Status),
		"mediaBase":      "http://localhost:9000/test",
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

// RegisterAPI godoc
// @title battery life API
// @version 1.0
// @description battery life — API для управления заявками расчёта времени работы аккумуляторов
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// RegisterAPI регистрирует маршруты REST API под префиксом /api
func (h *Handler) RegisterAPI(router *gin.Engine) {
	api := router.Group("/api")

	public := api.Group("/")
	{
		public.GET("/battery_life_types", h.APIGetBatteryTypes)
		public.GET("/battery_life_type/:id", h.APIGetBatteryType)
		public.POST("/users/signup", h.APICreateUser)
		public.POST("/users/signin", h.APISignIn)
	}

	optional := api.Group("/")
	optional.Use(h.OptionalAuthMiddleware())
	{
		optional.GET("/battery_life/battery_life-cart", h.APIGetBatteryLifeCart)
	}

	authorized := api.Group("/")
	authorized.Use(h.AuthMiddleware(false))
	{
		authorized.POST("/battery_life_type/create-battery_life_type", h.APICreateBatteryType)
		authorized.GET("/battery_life/all-battery_life", h.APIGetBatteryLives)
		authorized.GET("/battery_life/:id", h.APIGetBatteryLife)
		authorized.PUT("/battery_life/:id/edit-battery_life", h.APIEditBatteryLife)
		authorized.PUT("/battery_life/:id/form-battery_life", h.APIFormBatteryLife)
		authorized.DELETE("/battery_life/:id/delete-battery_life", h.APIDeleteBatteryLife)
		authorized.POST("/battery_life_item/add/:battery_life_type_id", h.APIAddToBatteryLife)
		authorized.DELETE("/battery_life_item/:battery_life_type_id/:battery_life_id", h.APIDeleteFromBatteryLife)
		authorized.PUT("/battery_life_item/:battery_life_type_id/:battery_life_id", h.APIEditInBatteryLife)
		authorized.POST("/users/signout", h.APISignOut)
	}

	moderator := api.Group("/")
	moderator.Use(h.AuthMiddleware(true))
	{
		moderator.PUT("/battery_life/:id/finish-battery_life", h.APIFinishBatteryLife)
	}

	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}
