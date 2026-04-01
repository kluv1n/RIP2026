package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"RIP2026/internal/app/repository"
	"RIP2026/internal/app/serializer"
	"github.com/gin-gonic/gin"
)

// APIGetBatteryLifeCart godoc
// @Summary Получить корзину заявки
// @Description Возвращает информацию о текущем черновике пользователя или статус `no_draft`.
// @Tags battery_lives
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /battery_life/battery_life-cart [get]
func (h *Handler) APIGetBatteryLifeCart(ctx *gin.Context) {
	creatorID, err := currentUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":           "no_draft",
			"count":            0,
			"strategies_count": 0,
		})
		return
	}
	count := h.Repository.GetCartItemCount(creatorID)
	if count == 0 {
		load, err := h.Repository.CheckCurrentDraft(creatorID)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":           "no_draft",
				"count":            0,
				"strategies_count": 0,
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

// APIGetBatteryLives godoc
// @Summary Получить список заявок
// @Description Возвращает заявки пользователя, а для модератора - все заявки. Поддерживает фильтрацию по датам и статусу.
// @Tags battery_lives
// @Produce json
// @Param from-date query string false "Начальная дата (YYYY-MM-DD)"
// @Param to-date query string false "Конечная дата (YYYY-MM-DD)"
// @Param status query string false "Статус заявки"
// @Success 200 {array} serializer.BatteryLifeJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life/all-battery_lives [get]
func (h *Handler) APIGetBatteryLives(ctx *gin.Context) {
	userID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(userID))

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

// APIGetBatteryLife godoc
// @Summary Получить заявку по ID
// @Description Возвращает полную информацию о заявке и ее позициях.
// @Tags battery_lives
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life/{id} [get]
func (h *Handler) APIGetBatteryLife(ctx *gin.Context) {
	userID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(userID))

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
		"items":        itemsResp,
	})
}

// APIEditBatteryLife godoc
// @Summary Изменить заявку
// @Description Обновляет данные черновика заявки.
// @Tags battery_lives
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param battery_life body serializer.BatteryLifeJSON true "Новые данные заявки"
// @Success 200 {object} serializer.BatteryLifeJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life/{id}/edit-battery_life [put]
func (h *Handler) APIEditBatteryLife(ctx *gin.Context) {
	userID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(userID))

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

// APIFormBatteryLife godoc
// @Summary Сформировать заявку
// @Description Переводит черновик в статус `formed`.
// @Tags battery_lives
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} serializer.BatteryLifeJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life/{id}/form-battery_life [put]
func (h *Handler) APIFormBatteryLife(ctx *gin.Context) {
	userID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(userID))

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

// APIFinishBatteryLife godoc
// @Summary Завершить заявку
// @Description Изменяет статус заявки на `completed` или `rejected`. Доступно только модератору.
// @Tags battery_lives
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param status body serializer.StatusJSON true "Новый статус заявки"
// @Success 200 {object} serializer.BatteryLifeJSON
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life/{id}/finish-battery_life [put]
func (h *Handler) APIFinishBatteryLife(ctx *gin.Context) {
	userID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(userID))

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

// APIDeleteBatteryLife godoc
// @Summary Удалить заявку
// @Description Выполняет логическое удаление черновика заявки.
// @Tags battery_lives
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /battery_life/{id}/delete-battery_life [delete]
func (h *Handler) APIDeleteBatteryLife(ctx *gin.Context) {
	userID, err := currentUserID(ctx)
	if err != nil {
		h.apiError(ctx, http.StatusUnauthorized, repository.ErrNotAllowed)
		return
	}
	h.Repository.SetUserID(int(userID))

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
