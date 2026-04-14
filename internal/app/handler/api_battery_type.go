package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"RIP2026/internal/app/repository"
	"RIP2026/internal/app/serializer"
	"github.com/gin-gonic/gin"
)

// APIGetBatteryTypes godoc
// @Summary Получить список аккумуляторов
// @Description Возвращает все аккумуляторы или фильтрует по названию.
// @Tags battery_types
// @Produce json
// @Param title query string false "Название аккумулятора для поиска"
// @Success 200 {object} serializer.BatteryTypeListJSON
// @Failure 500 {object} map[string]string
// @Router /battery_types [get]
func (h *Handler) APIGetBatteryTypes(ctx *gin.Context) {
	title := ctx.Query("title")
	if title == "" {
		types, err := h.Repository.GetBatteryTypesAPI()
		if err != nil {
			h.apiError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, serializer.BatteryTypesToListJSON(types))
		return
	}
	types, err := h.Repository.GetBatteryTypesByTitleAPI(title)
	if err != nil {
		h.apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, serializer.BatteryTypesToListJSON(types))
}

// APIGetBatteryType godoc
// @Summary Получить аккумулятор по ID
// @Description Возвращает данные одного аккумулятора.
// @Tags battery_types
// @Produce json
// @Param id path int true "ID аккумулятора"
// @Success 200 {object} serializer.BatteryTypeJSON
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /battery_type/{id} [get]
func (h *Handler) APIGetBatteryType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	bt, err := h.Repository.GetBatteryTypeAPI(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.apiError(ctx, http.StatusNotFound, err)
		} else {
			h.apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.BatteryTypeToJSON(*bt))
}

// APICreateBatteryType godoc
// @Summary Создать аккумулятор
// @Description Создает новый тип аккумулятора. Доступно авторизованному пользователю.
// @Tags battery_types
// @Accept json
// @Produce json
// @Param battery_type body serializer.BatteryTypeJSON true "Данные аккумулятора"
// @Success 201 {object} serializer.BatteryTypeJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_type/create-battery_type [post]
func (h *Handler) APICreateBatteryType(ctx *gin.Context) {
	contentType := ctx.GetHeader("Content-Type")
	var j serializer.BatteryTypeJSON
	if strings.HasPrefix(contentType, "application/json") {
		if err := ctx.BindJSON(&j); err != nil {
			h.apiError(ctx, http.StatusBadRequest, err)
			return
		}
	} else {
		title := ctx.PostForm("title")
		shortDesc := ctx.PostForm("short_description")
		desc := ctx.PostForm("description")
		if title == "" || desc == "" {
			h.apiError(ctx, http.StatusBadRequest, fmt.Errorf("title and description are required"))
			return
		}
		capacityMah, _ := strconv.Atoi(ctx.PostForm("capacity_mah"))
		if capacityMah <= 0 {
			capacityMah = 1000
		}
		voltageV := 3.7
		if v := ctx.PostForm("voltage_v"); v != "" {
			fmt.Sscanf(v, "%f", &voltageV)
		}
		j = serializer.BatteryTypeJSON{
			Title:            title,
			CapacityMah:      capacityMah,
			VoltageV:         voltageV,
			ShortDescription: shortDesc,
			Description:      desc,
		}
	}
	bt, err := h.Repository.CreateBatteryTypeAPI(j)
	if err != nil {
		h.apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	if imageFile, err := ctx.FormFile("image"); err == nil {
		updated, err := h.Repository.AddPhoto(ctx, int(bt.ID), imageFile)
		if err != nil {
			h.apiError(ctx, http.StatusInternalServerError, err)
			return
		}
		bt = *updated
	}
	if videoFile, err := ctx.FormFile("video"); err == nil {
		updated, err := h.Repository.AddVideo(ctx, int(bt.ID), videoFile)
		if err != nil {
			h.apiError(ctx, http.StatusInternalServerError, err)
			return
		}
		bt = *updated
	}
	ctx.Header("Location", fmt.Sprintf("/api/battery_types/%d", bt.ID))
	ctx.JSON(http.StatusCreated, serializer.BatteryTypeToJSON(bt))
}
