package models

// BatteryType — услуга (тип аккумулятора): наименование, описание, статус удалён/действует, url изображения nullable
type BatteryType struct {
	ID               uint   `gorm:"primaryKey"`
	Title            string `gorm:"type:varchar(200);not null"`
	CapacityMah      int    `gorm:"not null"`
	VoltageV         float64 `gorm:"type:decimal(5,2);not null"`
	Photo            *string `gorm:"type:varchar(255)"` // nullable
	Video            string  `gorm:"type:varchar(255)"`
	ShortDescription string  `gorm:"type:text"`
	Description      string  `gorm:"type:text"`
	IsDeleted        bool    `gorm:"type:boolean;default:false;not null"`
}

func (BatteryType) TableName() string { return "battery_types" }
