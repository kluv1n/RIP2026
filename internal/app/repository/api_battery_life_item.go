package repository

import (
	"errors"
	"fmt"

	"RIP2026/internal/app/models"
	"RIP2026/internal/app/serializer"
	"gorm.io/gorm"
)

func (r *Repository) AddBatteryToBatteryLifeAPI(creatorID uint, batteryTypeID int, currentMa, quantity int) (models.BatteryLife, bool, error) {
	load, created, err := r.GetBatteryLifeDraft(creatorID)
	if err != nil {
		return models.BatteryLife{}, false, err
	}
	var bt models.BatteryType
	if err := r.db.Where("id = ? AND is_deleted = ?", batteryTypeID, false).First(&bt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BatteryLife{}, false, fmt.Errorf("%w: тип аккумулятора с id %d", ErrNotFound, batteryTypeID)
		}
		return models.BatteryLife{}, false, err
	}
	var count int64
	r.db.Model(&models.BatteryLifeItem{}).
		Where("battery_life_id = ? AND battery_type_id = ?", load.ID, batteryTypeID).
		Count(&count)
	if count > 0 {
		return models.BatteryLife{}, false, fmt.Errorf("%w: услуга уже в заявке", ErrAlreadyExists)
	}
	runtimeHours := RuntimeHours(bt.CapacityMah, currentMa)
	if quantity <= 0 {
		quantity = 1
	}
	if currentMa <= 0 {
		currentMa = 100
	}
	item := models.BatteryLifeItem{
		BatteryLifeID: load.ID,
		BatteryTypeID: uint(batteryTypeID),
		CurrentMa:     currentMa,
		Quantity:      quantity,
		RuntimeHours:  runtimeHours,
		SortOrder:     0,
	}
	if err := r.db.Create(&item).Error; err != nil {
		return models.BatteryLife{}, false, err
	}
	return load, created, nil
}

func (r *Repository) DeleteFromBatteryLife(lifeID, batteryTypeID int) (models.BatteryLife, error) {
	var load models.BatteryLife
	if err := r.db.Where("id = ?", lifeID).First(&load).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BatteryLife{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, lifeID)
		}
		return models.BatteryLife{}, err
	}
	if load.Status != models.StatusDraft {
		return models.BatteryLife{}, fmt.Errorf("%w: можно удалять только из черновика", ErrNotAllowed)
	}
	if err := r.db.Where("battery_life_id = ? AND battery_type_id = ?", lifeID, batteryTypeID).
		Delete(&models.BatteryLifeItem{}).Error; err != nil {
		return models.BatteryLife{}, err
	}
	return load, nil
}

func (r *Repository) EditInBatteryLife(lifeID, batteryTypeID int, j serializer.BatteryLifeItemJSON) (models.BatteryLifeItem, error) {
	var item models.BatteryLifeItem
	if err := r.db.Where("battery_life_id = ? AND battery_type_id = ?", lifeID, batteryTypeID).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BatteryLifeItem{}, fmt.Errorf("%w: позиция в заявке", ErrNotFound)
		}
		return models.BatteryLifeItem{}, err
	}
	var load models.BatteryLife
	if err := r.db.Where("id = ?", lifeID).First(&load).Error; err != nil {
		return models.BatteryLifeItem{}, err
	}
	if load.Status != models.StatusDraft {
		return models.BatteryLifeItem{}, fmt.Errorf("%w: можно редактировать только черновик", ErrNotAllowed)
	}
	updates := map[string]interface{}{
		"quantity":    j.Quantity,
		"current_ma":  j.CurrentMa,
		"sort_order": j.SortOrder,
	}
	var bt models.BatteryType
	if err := r.db.First(&bt, item.BatteryTypeID).Error; err == nil {
		rh := RuntimeHours(bt.CapacityMah, j.CurrentMa)
		updates["runtime_hours"] = rh
	}
	if err := r.db.Model(&item).Updates(updates).Error; err != nil {
		return models.BatteryLifeItem{}, err
	}
	r.db.Where("battery_life_id = ? AND battery_type_id = ?", lifeID, batteryTypeID).
		Preload("BatteryType").First(&item)
	return item, nil
}
