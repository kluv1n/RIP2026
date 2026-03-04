package models

// BatteryLifeItem — м-м заявки-услуги: составной уникальный ключ (battery_life_id, battery_type_id), количество, ток, время работы
type BatteryLifeItem struct {
	ID             uint    `gorm:"primaryKey"`
	BatteryLifeID  uint    `gorm:"not null;uniqueIndex:idx_battery_life_type"`
	BatteryTypeID  uint    `gorm:"not null;uniqueIndex:idx_battery_life_type"`
	CurrentMa      int     `gorm:"not null"`  // потребляемый ток, мА
	Quantity       int     `gorm:"not null"`  // м-м: количество
	RuntimeHours   float64 `gorm:"type:decimal(12,4)"` // время работы (ёмкость/ток), ч

	BatteryLife BatteryLife `gorm:"foreignKey:BatteryLifeID"`
	BatteryType BatteryType `gorm:"foreignKey:BatteryTypeID"`
}

func (BatteryLifeItem) TableName() string { return "battery_life_items" }
