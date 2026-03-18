package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"RIP2026/internal/app/repository"
	"RIP2026/internal/app/serializer"
)

func (h *Handler) APIGetBatteryLifeCart(ctx *gin.Context) {
	creatorID := uint(h.Repository.GetCreatorID())
	count := h.Repository.GetCartItemCount(creatorID)
	if count == 0 {
		load, err := h.Repository.CheckCurrentDraft(creatorID)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":            "no_draft",
				"count":             0,
				"strategies_count":  0,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"id":               load.ID,
			"count":            0,
			"strategies_count": 0,
		})
		return
	}
	loadID := h.Repository.GetActiveDraftID(creatorID)
	ctx.JSON(http.StatusOK, gin.H{
		"id":               loadID,
		"count":            count,
		"strategies_count": count,
	})
}

func (h *Handler) APIGetBatteryLives(ctx *gin.Context) {
	fromDate := ctx.Query("from_date")
	if fromDate == "" {
		fromDate = ctx.Query("from-date")
	}
	var from, to time.Time
	if fromDate != "" {
		t, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.apiError(ctx, http.StatusBadRequest, err)
			return
		}
		from = t
	}
	toDate := ctx.Query("to_date")
	if toDate == "" {
		toDate = ctx.Query("to-date")
	}
	if toDate != "" {
		t, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.apiError(ctx, http.StatusBadRequest, err)
			return
		}
		to = t
	}
	status := ctx.Query("status")
	loads, err := h.Repository.GetAllBatteryLives(from, to, status)
	if err != nil {
		h.apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]serializer.BatteryLifeJSON, 0, len(loads))
	for _, load := range loads {
		creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
		completedCount, _ := h.Repository.GetCompletedItemCount(load.ID)
		resp = append(resp, serializer.BatteryLifeToJSON(load, creatorLogin, moderatorLogin, completedCount))
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) APIGetBatteryLife(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.GetSingleBatteryLife(id)
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
	items, err := h.Repository.GetBatteryLifeItems(id)
	if err != nil {
		h.apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
	completedCount, _ := h.Repository.GetCompletedItemCount(load.ID)
	itemsResp := make([]serializer.BatteryLifeItemDetailJSON, 0, len(items))
	for _, item := range items {
		itemsResp = append(itemsResp, serializer.BatteryLifeItemDetailToJSON(item))
	}
	ctx.JSON(http.StatusOK, gin.H{
		"battery_life": serializer.BatteryLifeToJSON(load, creatorLogin, moderatorLogin, completedCount),
		"items":       itemsResp,
	})
}

func (h *Handler) APIEditBatteryLife(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	var j serializer.BatteryLifeJSON
	if err := ctx.BindJSON(&j); err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.EditBatteryLife(id, j)
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

func (h *Handler) APIFormBatteryLife(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.FormBatteryLife(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.apiError(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.apiError(ctx, http.StatusForbidden, err)
		} else {
			h.apiError(ctx, http.StatusBadRequest, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
	completedCount, _ := h.Repository.GetCompletedItemCount(load.ID)
	ctx.JSON(http.StatusOK, serializer.BatteryLifeToJSON(load, creatorLogin, moderatorLogin, completedCount))
}

func (h *Handler) APIFinishBatteryLife(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	var statusJ serializer.StatusJSON
	if err := ctx.BindJSON(&statusJ); err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	load, err := h.Repository.FinishBatteryLife(id, statusJ.Status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.apiError(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.apiError(ctx, http.StatusForbidden, err)
		} else {
			h.apiError(ctx, http.StatusBadRequest, err)
		}
		return
	}
	creatorLogin, moderatorLogin, _ := h.Repository.GetModeratorAndCreatorLogin(load)
	completedCount, _ := h.Repository.GetCompletedItemCount(load.ID)
	ctx.JSON(http.StatusOK, serializer.BatteryLifeToJSON(load, creatorLogin, moderatorLogin, completedCount))
}

func (h *Handler) APIDeleteBatteryLife(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.apiError(ctx, http.StatusBadRequest, err)
		return
	}
	_, err = h.Repository.DeleteBatteryLifeAPI(id)
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
	ctx.JSON(http.StatusOK, gin.H{"message": "Заявка удалена"})
}
