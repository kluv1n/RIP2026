package serializer

import (
	"os"
	"strings"

	"RIP2026/internal/app/models"
)

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

type BatteryTypeListJSON struct {
	IDs   []uint            `json:"ids"`
	Items []BatteryTypeJSON `json:"items"`
}

func BatteryTypeToJSON(b models.BatteryType) BatteryTypeJSON {
	photo := ""
	if b.Photo != nil {
		photo = mediaURL(*b.Photo)
	}
	return BatteryTypeJSON{
		ID:               b.ID,
		Title:            b.Title,
		CapacityMah:      b.CapacityMah,
		VoltageV:         b.VoltageV,
		Photo:            photo,
		Video:            mediaURL(b.Video),
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

func BatteryTypesToListJSON(items []models.BatteryType) BatteryTypeListJSON {
	resp := BatteryTypeListJSON{
		IDs:   make([]uint, 0, len(items)),
		Items: make([]BatteryTypeJSON, 0, len(items)),
	}
	for _, item := range items {
		resp.IDs = append(resp.IDs, item.ID)
		resp.Items = append(resp.Items, BatteryTypeToJSON(item))
	}
	return resp
}

func mediaURL(objectName string) string {
	if objectName == "" {
		return ""
	}
	if strings.HasPrefix(objectName, "http://") || strings.HasPrefix(objectName, "https://") {
		return objectName
	}
	base := os.Getenv("MINIO_URL")
	if base == "" {
		base = "http://localhost:9000/batteries"
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(objectName, "/")
}
