package models

import "time"

// Статусы заявки (API и БД — на английском, как в методичке/скринах)
const (
	StatusDraft     = "draft"
	StatusDeleted   = "deleted"
	StatusFormed    = "formed"
	StatusCompleted = "completed"
	StatusRejected  = "rejected"
)

// BatteryLife — заявка (расчёт времени работы): id, статус, дата создания, создатель; даты формирования/завершения, модератор
type BatteryLife struct {
	ID                uint       `gorm:"primaryKey"`
	Status            string     `gorm:"type:varchar(20);not null"`
	CreatedAt         time.Time  `gorm:"not null"`
	CreatorID         uint       `gorm:"not null"`
	FormedAt          *time.Time // дата формирования (2 действия создателя)
	CompletedAt       *time.Time // дата завершения (2 действия модератора)
	ModeratorID       *uint
	Title             string  `gorm:"type:varchar(255)"`
	Description       string  `gorm:"type:text"`
	TotalRuntimeHours float64 `gorm:"type:decimal(12,4)"` // рассчитывается при завершении

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator *User `gorm:"foreignKey:ModeratorID"`
	Items     []BatteryLifeItem `gorm:"foreignKey:BatteryLifeID"`
}

func (BatteryLife) TableName() string { return "battery_lives" }
