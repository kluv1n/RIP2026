package serializer

import "RIP2026/internal/app/models"

type BatteryTypeJSON struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	CapacityMah      int     `json:"capacity_mah"`
	VoltageV         float64 `json:"voltage_v"`
	Photo            string  `json:"photo"`
	Video            string  `json:"video"`
	ShortDescription string  `json:"short_description"`
	Description      string  `json:"description"`
	IsDeleted        bool    `json:"is_deleted"`
}

func BatteryTypeToJSON(b models.BatteryType) BatteryTypeJSON {
	photo := ""
	if b.Photo != nil {
		photo = *b.Photo
	}
	return BatteryTypeJSON{
		ID:               b.ID,
		Title:            b.Title,
		CapacityMah:      b.CapacityMah,
		VoltageV:         b.VoltageV,
		Photo:            photo,
		Video:            b.Video,
		ShortDescription: b.ShortDescription,
		Description:      b.Description,
		IsDeleted:        b.IsDeleted,
	}
}

func BatteryTypeFromJSON(j BatteryTypeJSON) models.BatteryType {
	return models.BatteryType{
		Title:            j.Title,
		CapacityMah:      j.CapacityMah,
		VoltageV:         j.VoltageV,
		ShortDescription: j.ShortDescription,
		Description:      j.Description,
	}
}
