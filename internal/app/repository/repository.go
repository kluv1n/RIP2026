package repository

import (
	"fmt"
	"time"

	"RIP2026/internal/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultCreatorID = 1

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// DTO для шаблонов (совместимы с прежними полями)
type BatteryType struct {
	ID               int
	Title            string
	CapacityMah      int
	VoltageV         float64
	Photo            string
	Video            string
	ShortDescription string
	Description      string
}

type BatteryLifeItem struct {
	Battery      BatteryType
	CurrentMa    int
	Mm           string
	Quantity     int
	RuntimeHours float64
	RuntimeTotal float64
}

type BatteryLife struct {
	ID                int
	Title             string
	Description       string
	Items             []BatteryLifeItem
	ItemCount         int
	TotalRuntimeHours float64
	Status            string
}

func RuntimeHours(capacityMah, currentMa int) float64 {
	if currentMa <= 0 {
		return 0
	}
	return float64(capacityMah) / float64(currentMa)
}

func (r *Repository) GetBatteryTypes() ([]BatteryType, error) {
	var rows []models.BatteryType
	if err := r.db.Where("is_deleted = ?", false).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]BatteryType, len(rows))
	for i := range rows {
		photo := ""
		if rows[i].Photo != nil {
			photo = *rows[i].Photo
		}
		out[i] = BatteryType{
			ID:               int(rows[i].ID),
			Title:            rows[i].Title,
			CapacityMah:      rows[i].CapacityMah,
			VoltageV:         rows[i].VoltageV,
			Photo:            photo,
			Video:            rows[i].Video,
			ShortDescription: rows[i].ShortDescription,
			Description:      rows[i].Description,
		}
	}
	return out, nil
}

func (r *Repository) GetBattery(id int) (BatteryType, error) {
	var m models.BatteryType
	if err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&m).Error; err != nil {
		return BatteryType{}, err
	}
	photo := ""
	if m.Photo != nil {
		photo = *m.Photo
	}
	return BatteryType{
		ID:               int(m.ID),
		Title:            m.Title,
		CapacityMah:      m.CapacityMah,
		VoltageV:         m.VoltageV,
		Photo:            photo,
		Video:            m.Video,
		ShortDescription: m.ShortDescription,
		Description:      m.Description,
	}, nil
}

func (r *Repository) GetBatteryByTitle(title string) ([]BatteryType, error) {
	var rows []models.BatteryType
	if err := r.db.Where("is_deleted = ? AND title ILIKE ?", false, "%"+title+"%").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]BatteryType, len(rows))
	for i := range rows {
		photo := ""
		if rows[i].Photo != nil {
			photo = *rows[i].Photo
		}
		out[i] = BatteryType{
			ID:               int(rows[i].ID),
			Title:            rows[i].Title,
			CapacityMah:      rows[i].CapacityMah,
			VoltageV:         rows[i].VoltageV,
			Photo:            photo,
			Video:            rows[i].Video,
			ShortDescription: rows[i].ShortDescription,
			Description:      rows[i].Description,
		}
	}
	return out, nil
}

func (r *Repository) GetDraftID(creatorID uint) (uint, bool) {
	var lives []models.BatteryLife
	r.db.Where("creator_id = ? AND status = ?", creatorID, models.StatusDraft).Select("id").Limit(1).Find(&lives)
	if len(lives) == 0 {
		return 0, false
	}
	return lives[0].ID, true
}

func (r *Repository) GetCartCount(creatorID uint) int64 {
	draftID, ok := r.GetDraftID(creatorID)
	if !ok {
		return 0
	}
	var count int64
	r.db.Model(&models.BatteryLifeItem{}).Where("battery_life_id = ?", draftID).Count(&count)
	return count
}

// GetBatteryLives возвращает заявки пользователя, кроме удалённых
func (r *Repository) GetBatteryLives(creatorID uint) ([]BatteryLife, error) {
	var rows []models.BatteryLife
	if err := r.db.Where("creator_id = ? AND status != ?", creatorID, models.StatusDeleted).
		Preload("Items.BatteryType").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]BatteryLife, 0, len(rows))
	for _, row := range rows {
		bl := rowToBatteryLife(&row)
		out = append(out, bl)
	}
	return out, nil
}

func (r *Repository) GetBatteryLife(id int, creatorID uint) (BatteryLife, error) {
	var row models.BatteryLife
	if err := r.db.Where("id = ? AND creator_id = ? AND status != ?", id, creatorID, models.StatusDeleted).
		Preload("Items.BatteryType").First(&row).Error; err != nil {
		return BatteryLife{}, err
	}
	return rowToBatteryLife(&row), nil
}

func rowToBatteryLife(row *models.BatteryLife) BatteryLife {
	items := make([]BatteryLifeItem, 0, len(row.Items))
	for i := range row.Items {
		bt := BatteryType{}
		if row.Items[i].BatteryType.ID != 0 {
			photo := ""
			if row.Items[i].BatteryType.Photo != nil {
				photo = *row.Items[i].BatteryType.Photo
			}
			bt = BatteryType{
				ID:               int(row.Items[i].BatteryType.ID),
				Title:            row.Items[i].BatteryType.Title,
				CapacityMah:      row.Items[i].BatteryType.CapacityMah,
				VoltageV:         row.Items[i].BatteryType.VoltageV,
				Photo:            photo,
				Video:            row.Items[i].BatteryType.Video,
				ShortDescription: row.Items[i].BatteryType.ShortDescription,
				Description:      row.Items[i].BatteryType.Description,
			}
		}
		qty := row.Items[i].Quantity
		if qty <= 0 {
			qty = 1
		}
		rh := row.Items[i].RuntimeHours
		rt := rh * float64(qty)
		items = append(items, BatteryLifeItem{
			Battery:      bt,
			CurrentMa:    row.Items[i].CurrentMa,
			Mm:           fmt.Sprintf("%d", qty),
			Quantity:     qty,
			RuntimeHours: rh,
			RuntimeTotal: rt,
		})
	}
	total := 0.0
	for _, it := range items {
		total += it.RuntimeTotal
	}
	return BatteryLife{
		ID:                int(row.ID),
		Title:             row.Title,
		Description:       row.Description,
		Items:             items,
		ItemCount:         len(items),
		TotalRuntimeHours: total,
		Status:            row.Status,
	}
}

func (r *Repository) GetBatteryLifeForBattery(batteryID int, creatorID uint) (*BatteryLifeItem, error) {
	draftID, ok := r.GetDraftID(creatorID)
	if !ok {
		return nil, nil
	}
	var items []models.BatteryLifeItem
	r.db.Where("battery_life_id = ? AND battery_type_id = ?", draftID, batteryID).
		Preload("BatteryType").Limit(1).Find(&items)
	if len(items) == 0 {
		return nil, nil
	}
	item := items[0]
	photo := ""
	if item.BatteryType.Photo != nil {
		photo = *item.BatteryType.Photo
	}
	qty := item.Quantity
	if qty <= 0 {
		qty = 1
	}
	rh := item.RuntimeHours
	return &BatteryLifeItem{
		Battery: BatteryType{
			ID:               int(item.BatteryType.ID),
			Title:            item.BatteryType.Title,
			CapacityMah:      item.BatteryType.CapacityMah,
			VoltageV:         item.BatteryType.VoltageV,
			Photo:            photo,
			Video:            item.BatteryType.Video,
			ShortDescription: item.BatteryType.ShortDescription,
			Description:      item.BatteryType.Description,
		},
		CurrentMa:    item.CurrentMa,
		Mm:           fmt.Sprintf("%d", qty),
		Quantity:     qty,
		RuntimeHours: rh,
		RuntimeTotal: rh * float64(qty),
	}, nil
}

// AddBatteryToBatteryLife добавляет услугу в заявку (черновик). Создаёт черновик, если его нет. Через ORM.
func (r *Repository) AddBatteryToBatteryLife(creatorID uint, batteryTypeID int, currentMa, quantity int) error {
	var draft models.BatteryLife
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, models.StatusDraft).First(&draft).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			draft = models.BatteryLife{
				Status:    models.StatusDraft,
				CreatedAt: time.Now(),
				CreatorID: creatorID,
				Title:     "Расчёт времени работы",
				Description: "Время (ч) = ёмкость (мА·ч) / ток (мА).",
			}
			if err = r.db.Create(&draft).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	var bt models.BatteryType
	if err := r.db.Where("id = ? AND is_deleted = ?", batteryTypeID, false).First(&bt).Error; err != nil {
		return err
	}
	runtimeHours := RuntimeHours(bt.CapacityMah, currentMa)
	if quantity <= 0 {
		quantity = 1
	}

	var existing models.BatteryLifeItem
	err = r.db.Where("battery_life_id = ? AND battery_type_id = ?", draft.ID, batteryTypeID).First(&existing).Error
	if err == nil {
		// уже в заявке — обновляем ток, добавляем к количеству
		newQty := existing.Quantity + quantity
		newMa := currentMa
		newRuntime := RuntimeHours(bt.CapacityMah, newMa)
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"current_ma":     newMa,
			"quantity":      newQty,
			"runtime_hours": newRuntime,
		}).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	item := models.BatteryLifeItem{
		BatteryLifeID: draft.ID,
		BatteryTypeID: uint(batteryTypeID),
		CurrentMa:     currentMa,
		Quantity:      quantity,
		RuntimeHours:  runtimeHours,
	}
	return r.db.Create(&item).Error
}

// DeleteBatteryLife логическое удаление заявки: SQL UPDATE без ORM
func (r *Repository) DeleteBatteryLife(id int, creatorID uint) error {
	res := r.db.Exec(
		"UPDATE battery_lives SET status = ? WHERE id = ? AND creator_id = ? AND status = ?",
		models.StatusDeleted, id, creatorID, models.StatusDraft,
	)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("заявка не найдена или уже удалена")
	}
	return nil
}
