package serializer

import (
	"time"

	"RIP2026/internal/app/models"
)

type BatteryLifeJSON struct {
	ID                 uint        `json:"id"`
	Status             string      `json:"status"`
	CreatedAt          time.Time   `json:"created_at"`
	CreatorLogin       string      `json:"creator_login"`
	ModeratorLogin     *string     `json:"moderator_login,omitempty"` // опционально, как в SystemLoadToJSON
	FormedAt           *time.Time  `json:"formed_at"`
	CompletedAt        *time.Time  `json:"completed_at"`
	Title              string      `json:"title"`
	Description        string      `json:"description"`
	TotalRuntimeHours  float64     `json:"total_runtime_hours"`
	CompletedItemCount int         `json:"completed_item_count"`
}

func BatteryLifeToJSON(load models.BatteryLife, creatorLogin, moderatorLogin string, completedItemCount int) BatteryLifeJSON {
	var mLogin *string
	if moderatorLogin != "" {
		mLogin = &moderatorLogin
	}
	var formedAt *time.Time
	if load.FormedAt != nil {
		formedAt = load.FormedAt
	}
	var completedAt *time.Time
	if load.CompletedAt != nil {
		completedAt = load.CompletedAt
	}
	return BatteryLifeJSON{
		ID:                 load.ID,
		Status:             load.Status,
		CreatedAt:          load.CreatedAt,
		CreatorLogin:       creatorLogin,
		ModeratorLogin:     mLogin,
		FormedAt:           formedAt,
		CompletedAt:        completedAt,
		Title:              load.Title,
		Description:        load.Description,
		TotalRuntimeHours:  load.TotalRuntimeHours,
		CompletedItemCount: completedItemCount,
	}
}

type StatusJSON struct {
	Status string `json:"status"`
}
