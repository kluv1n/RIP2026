package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"RIP2026/internal/app/repository"
	"RIP2026/internal/app/serializer"
)

func (h *Handler) APIAddToBatteryLife(ctx *gin.Context) {
	batteryTypeIDStr := ctx.Param("battery_type_id")
	batteryTypeID, err := strconv.Atoi(batteryTypeIDStr)
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	creatorID := uint(h.Repository.GetCreatorID())
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
		ctx.Header("Location", fmt.Sprintf("/api/battery_lives/%d", load.ID))
		status = http.StatusCreated
	}
	ctx.JSON(status, serializer.BatteryLifeToJSON(load, creatorLogin, moderatorLogin, completedCount))
}

func (h *Handler) APIDeleteFromBatteryLife(ctx *gin.Context) {
	batteryTypeID, err := strconv.Atoi(ctx.Param("battery_type_id"))
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

func (h *Handler) APIEditInBatteryLife(ctx *gin.Context) {
	batteryTypeID, err := strconv.Atoi(ctx.Param("battery_type_id"))
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
