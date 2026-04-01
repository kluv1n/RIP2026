package repository

import (
	"errors"
	"fmt"
	"time"

	"RIP2026/internal/app/models"
	"RIP2026/internal/app/serializer"
	"gorm.io/gorm"
)

func (r *Repository) GetCartItemCount(creatorID uint) int64 {
	draftID, ok := r.GetDraftID(creatorID)
	if !ok {
		return 0
	}
	var count int64
	r.db.Model(&models.BatteryLifeItem{}).Where("battery_life_id = ?", draftID).Count(&count)
	return count
}

func (r *Repository) GetActiveDraftID(creatorID uint) uint {
	id, ok := r.GetDraftID(creatorID)
	if !ok {
		return 0
	}
	return id
}

func (r *Repository) CheckCurrentDraft(creatorID uint) (models.BatteryLife, error) {
	if creatorID == 0 {
		return models.BatteryLife{}, ErrNotAllowed
	}
	var load models.BatteryLife
	res := r.db.Where("creator_id = ? AND status = ?", creatorID, models.StatusDraft).Limit(1).Find(&load)
	if res.Error != nil {
		return models.BatteryLife{}, res.Error
	}
	if res.RowsAffected == 0 {
		return models.BatteryLife{}, ErrNoDraft
	}
	return load, nil
}

func (r *Repository) GetBatteryLifeDraft(creatorID uint) (models.BatteryLife, bool, error) {
	load, err := r.CheckCurrentDraft(creatorID)
	if err == ErrNoDraft {
		load = models.BatteryLife{
			Status:      models.StatusDraft,
			CreatedAt:   time.Now(),
			CreatorID:   creatorID,
			Title:       "Расчёт времени работы",
			Description: "Время (ч) = ёмкость (мА·ч) / ток (мА).",
		}
		if err := r.db.Create(&load).Error; err != nil {
			return models.BatteryLife{}, false, err
		}
		return load, true, nil
	}
	if err != nil {
		return models.BatteryLife{}, false, err
	}
	return load, false, nil
}

func (r *Repository) GetModeratorAndCreatorLogin(load models.BatteryLife) (string, string, error) {
	var creator models.User
	if err := r.db.Where("id = ?", load.CreatorID).First(&creator).Error; err != nil {
		return "", "", err
	}
	var moderatorLogin string
	if load.ModeratorID != nil && *load.ModeratorID != 0 {
		var moderator models.User
		if err := r.db.Where("id = ?", *load.ModeratorID).First(&moderator).Error; err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}
	return creator.Login, moderatorLogin, nil
}

func (r *Repository) GetCompletedItemCount(lifeID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.BatteryLifeItem{}).
		Where("battery_life_id = ? AND runtime_hours IS NOT NULL AND runtime_hours > 0", lifeID).
		Count(&count).Error
	return int(count), err
}

func (r *Repository) GetAllBatteryLives(from, to time.Time, status string) ([]models.BatteryLife, error) {
	var loads []models.BatteryLife
	sub := r.db.Where("status != ? AND status != ?", models.StatusDeleted, models.StatusDraft)
	if userID := r.GetUserID(); userID > 0 {
		user, err := r.GetUserByID(userID)
		if err == nil && !user.IsModerator {
			sub = sub.Where("creator_id = ?", user.ID)
		}
	}
	if !from.IsZero() {
		sub = sub.Where("formed_at >= ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("formed_at < ?", to.Add(24*time.Hour))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("id").Find(&loads).Error
	return loads, err
}

func (r *Repository) GetSingleBatteryLife(id int) (models.BatteryLife, error) {
	if id <= 0 {
		return models.BatteryLife{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
	}
	var load models.BatteryLife
	err := r.db.Where("id = ?", id).First(&load).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BatteryLife{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return models.BatteryLife{}, err
	}
	if load.Status == models.StatusDeleted {
		return models.BatteryLife{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	if userID := r.GetUserID(); userID > 0 {
		user, err := r.GetUserByID(userID)
		if err == nil && !user.IsModerator && load.CreatorID != user.ID {
			return models.BatteryLife{}, fmt.Errorf("%w: доступ только к своим заявкам", ErrNotAllowed)
		}
	}
	return load, nil
}

func (r *Repository) GetBatteryLifeItems(lifeID int) ([]models.BatteryLifeItem, error) {
	var items []models.BatteryLifeItem
	err := r.db.Where("battery_life_id = ?", lifeID).
		Preload("BatteryType").
		Order("sort_order, id").
		Find(&items).Error
	return items, err
}

func (r *Repository) EditBatteryLife(id int, j serializer.BatteryLifeJSON) (models.BatteryLife, error) {
	var load models.BatteryLife
	err := r.db.Where("id = ? AND status != ?", id, models.StatusDeleted).First(&load).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BatteryLife{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return models.BatteryLife{}, err
	}
	if load.Status != models.StatusDraft {
		return models.BatteryLife{}, fmt.Errorf("%w: можно редактировать только черновик", ErrNotAllowed)
	}
	if userID := r.GetUserID(); userID > 0 {
		user, err := r.GetUserByID(userID)
		if err == nil && !user.IsModerator && load.CreatorID != user.ID {
			return models.BatteryLife{}, fmt.Errorf("%w: можно редактировать только свою заявку", ErrNotAllowed)
		}
	}
	updates := map[string]interface{}{
		"title":       j.Title,
		"description": j.Description,
	}
	if err := r.db.Model(&load).Updates(updates).Error; err != nil {
		return models.BatteryLife{}, err
	}
	r.db.Where("id = ?", id).First(&load)
	return load, nil
}

func (r *Repository) FormBatteryLife(id int) (models.BatteryLife, error) {
	load, err := r.GetSingleBatteryLife(id)
	if err != nil {
		return models.BatteryLife{}, err
	}
	if load.Status != models.StatusDraft {
		return models.BatteryLife{}, fmt.Errorf("%w: только черновик можно сформировать", ErrNotAllowed)
	}
	if userID := r.GetUserID(); userID > 0 {
		user, err := r.GetUserByID(userID)
		if err == nil && !user.IsModerator && load.CreatorID != user.ID {
			return models.BatteryLife{}, fmt.Errorf("%w: можно сформировать только свою заявку", ErrNotAllowed)
		}
	}
	items, err := r.GetBatteryLifeItems(int(load.ID))
	if err != nil {
		return models.BatteryLife{}, err
	}
	if len(items) == 0 {
		return models.BatteryLife{}, fmt.Errorf("нельзя сформировать пустую заявку")
	}
	var totalHours float64
	for _, item := range items {
		var bt models.BatteryType
		if err := r.db.First(&bt, item.BatteryTypeID).Error; err != nil {
			return models.BatteryLife{}, err
		}
		rh := RuntimeHours(bt.CapacityMah, item.CurrentMa)
		qty := item.Quantity
		if qty <= 0 {
			qty = 1
		}
		totalHours += rh * float64(qty)
		r.db.Model(&models.BatteryLifeItem{}).
			Where("battery_life_id = ? AND battery_type_id = ?", load.ID, item.BatteryTypeID).
			Update("runtime_hours", rh)
	}
	formedAt := time.Now()
	if err := r.db.Model(&load).Updates(map[string]interface{}{
		"status":              models.StatusFormed,
		"formed_at":           formedAt,
		"total_runtime_hours": totalHours,
	}).Error; err != nil {
		return models.BatteryLife{}, err
	}
	load.Status = models.StatusFormed
	load.FormedAt = &formedAt
	load.TotalRuntimeHours = totalHours
	return load, nil
}

func (r *Repository) FinishBatteryLife(id int, status string) (models.BatteryLife, error) {
	if status != models.StatusCompleted && status != models.StatusRejected {
		return models.BatteryLife{}, fmt.Errorf("неверный статус: допустимы %s или %s", models.StatusCompleted, models.StatusRejected)
	}
	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return models.BatteryLife{}, err
	}
	if !user.IsModerator {
		return models.BatteryLife{}, fmt.Errorf("%w: вы не модератор", ErrNotAllowed)
	}
	load, err := r.GetSingleBatteryLife(id)
	if err != nil {
		return models.BatteryLife{}, err
	}
	if load.Status != models.StatusFormed {
		return models.BatteryLife{}, fmt.Errorf("%w: завершить/отклонить можно только сформированную заявку", ErrNotAllowed)
	}
	completedAt := time.Now()
	if err := r.db.Model(&load).Updates(map[string]interface{}{
		"status":       status,
		"completed_at": completedAt,
		"moderator_id": user.ID,
	}).Error; err != nil {
		return models.BatteryLife{}, err
	}
	load.Status = status
	load.CompletedAt = &completedAt
	load.ModeratorID = &user.ID
	return load, nil
}

func (r *Repository) DeleteBatteryLifeAPI(id int) (models.BatteryLife, error) {
	load, err := r.GetSingleBatteryLife(id)
	if err != nil {
		return models.BatteryLife{}, err
	}
	if load.Status != models.StatusDraft {
		return models.BatteryLife{}, fmt.Errorf("%w: удалить можно только черновик", ErrNotAllowed)
	}
	if userID := r.GetUserID(); userID > 0 {
		user, err := r.GetUserByID(userID)
		if err == nil && !user.IsModerator && load.CreatorID != user.ID {
			return models.BatteryLife{}, fmt.Errorf("%w: вы не создатель этой заявки", ErrNotAllowed)
		}
	}
	formedAt := time.Now()
	if err := r.db.Model(&load).Updates(map[string]interface{}{
		"status":    models.StatusDeleted,
		"formed_at": formedAt,
	}).Error; err != nil {
		return models.BatteryLife{}, err
	}
	load.Status = models.StatusDeleted
	load.FormedAt = &formedAt
	return load, nil
}
