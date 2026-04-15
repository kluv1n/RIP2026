package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"RIP2026/internal/app/repository"
	"RIP2026/internal/app/serializer"
	"github.com/gin-gonic/gin"
)

// APIAddToBatteryLife godoc
// @Summary Добавить аккумулятор в заявку
// @Description Добавляет аккумулятор в черновик пользователя. При отсутствии черновика он создается.
// @Tags battery life
// @Accept json
// @Produce json
// @Param battery_life_type_id path int true "ID типа аккумулятора"
// @Param item body serializer.BatteryLifeItemJSON false "Параметры позиции"
// @Success 200 {object} serializer.BatteryLifeJSON
// @Success 201 {object} serializer.BatteryLifeJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life_item/add/{battery_life_type_id} [post]
func (h *Handler) APIAddToBatteryLife(ctx *gin.Context) {
	creatorID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(creatorID))

	batteryTypeIDStr := ctx.Param("battery_life_type_id")
	batteryTypeID, err := strconv.Atoi(batteryTypeIDStr)
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	currentMa := 100
	if v := ctx.PostForm("current_ma"); v != "" {
		fmt.Sscanf(v, "%d", &currentMa)
	}
	if currentMa <= 0 {
		currentMa = 100
	}
	quantity := 1
	if v := ctx.PostForm("quantity"); v != "" {
		fmt.Sscanf(v, "%d", &quantity)
	}
	if quantity <= 0 {
		quantity = 1
	}
	var j serializer.BatteryLifeItemJSON
	_ = ctx.BindJSON(&j)
	if j.CurrentMa > 0 {
		currentMa = j.CurrentMa
	}
	if j.Quantity > 0 {
		quantity = j.Quantity
	}
	load, created, err := h.Repository.AddBatteryToBatteryLifeAPI(creatorID, batteryTypeID, currentMa, quantity)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.apiError(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.apiError(ctx, http.StatusConflict, err)
		} else {
			h.apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
	completedCount, _ := h.Repository.GetCompletedItemCount(load.ID)
	status := http.StatusOK
	if created {
		ctx.Header("Location", fmt.Sprintf("/api/battery_life/%d", load.ID))
		status = http.StatusCreated
	}
	ctx.JSON(status, serializer.BatteryLifeToJSON(load, creatorLogin, moderatorLogin, completedCount))
}

// APIDeleteFromBatteryLife godoc
// @Summary Удалить аккумулятор из заявки
// @Description Удаляет позицию из черновика заявки.
// @Tags battery life
// @Produce json
// @Param battery_life_type_id path int true "ID типа аккумулятора"
// @Param battery_life_id path int true "ID заявки"
// @Success 200 {object} serializer.BatteryLifeJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life_item/{battery_life_type_id}/{battery_life_id} [delete]
func (h *Handler) APIDeleteFromBatteryLife(ctx *gin.Context) {
	userID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(userID))

	batteryTypeID, err := strconv.Atoi(ctx.Param("battery_life_type_id"))
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	batteryLifeID, err := strconv.Atoi(ctx.Param("battery_life_id"))
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.DeleteFromBatteryLife(batteryLifeID, batteryTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.apiError(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.apiError(ctx, http.StatusForbidden, err)
		} else {
			h.apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
	completedCount, _ := h.Repository.GetCompletedItemCount(load.ID)
	ctx.JSON(http.StatusOK, serializer.BatteryLifeToJSON(load, creatorLogin, moderatorLogin, completedCount))
}

// APIEditInBatteryLife godoc
// @Summary Изменить позицию в заявке
// @Description Обновляет ток и количество для позиции в черновике заявки.
// @Tags battery life
// @Accept json
// @Produce json
// @Param battery_life_type_id path int true "ID типа аккумулятора"
// @Param battery_life_id path int true "ID заявки"
// @Param item body serializer.BatteryLifeItemJSON true "Новые данные позиции"
// @Success 200 {object} serializer.BatteryLifeItemJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life_item/{battery_life_type_id}/{battery_life_id} [put]
func (h *Handler) APIEditInBatteryLife(ctx *gin.Context) {
	userID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(userID))

	batteryTypeID, err := strconv.Atoi(ctx.Param("battery_life_type_id"))
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	batteryLifeID, err := strconv.Atoi(ctx.Param("battery_life_id"))
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.BatteryLifeItemJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	item, err := h.Repository.EditInBatteryLife(batteryLifeID, batteryTypeID, j)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.apiError(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.apiError(ctx, http.StatusForbidden, err)
		} else {
			h.apiError(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	ctx.JSON(http.StatusOK, serializer.BatteryLifeItemToJSON(item))
}
