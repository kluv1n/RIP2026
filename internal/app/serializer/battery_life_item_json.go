package serializer

import "RIP2026/internal/app/models"

type BatteryLifeItemJSON struct {
	BatteryLifeID uint    `json:"battery_life_id"`
	BatteryTypeID uint    `json:"battery_type_id"`
	CurrentMa     int     `json:"current_ma"`
	Quantity      int     `json:"quantity"`
	SortOrder     int     `json:"sort_order"`
	RuntimeHours  float64 `json:"runtime_hours"`
}

type BatteryLifeItemDetailJSON struct {
	BatteryLifeID uint            `json:"battery_life_id"`
	BatteryTypeID uint            `json:"battery_type_id"`
	CurrentMa     int             `json:"current_ma"`
	Quantity      int             `json:"quantity"`
	SortOrder     int             `json:"sort_order"`
	RuntimeHours  float64         `json:"runtime_hours"`
	BatteryType   BatteryTypeJSON `json:"battery_type"`
}

func BatteryLifeItemToJSON(item models.BatteryLifeItem) BatteryLifeItemJSON {
	return BatteryLifeItemJSON{
		BatteryLifeID: item.BatteryLifeID,
		BatteryTypeID: item.BatteryTypeID,
		CurrentMa:     item.CurrentMa,
		Quantity:      item.Quantity,
		SortOrder:     item.SortOrder,
		RuntimeHours:  item.RuntimeHours,
	}
}

func BatteryLifeItemDetailToJSON(item models.BatteryLifeItem) BatteryLifeItemDetailJSON {
	bt := BatteryTypeJSON{}
	if item.BatteryType.ID != 0 {
		bt = BatteryTypeToJSON(item.BatteryType)
	}
	return BatteryLifeItemDetailJSON{
		BatteryLifeID: item.BatteryLifeID,
		BatteryTypeID: item.BatteryTypeID,
		CurrentMa:     item.CurrentMa,
		Quantity:      item.Quantity,
		SortOrder:     item.SortOrder,
		RuntimeHours:  item.RuntimeHours,
		BatteryType:   bt,
	}
}
