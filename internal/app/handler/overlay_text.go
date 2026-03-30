package handler

import "strings"

// descriptionForVideoOverlay сокращает полное описание примерно в 2 раза (по словам),
// чтобы не было ни длинной простыни, ни одной короткой строки.
func descriptionForVideoOverlay(full string) string {
	full = strings.TrimSpace(full)
	if full == "" {
		return ""
	}
	words := strings.Fields(full)
	if len(words) <= 6 {
		return full
	}
	half := len(words) / 2
	if half < 3 {
		half = 3
	}
	return strings.Join(words[:half], " ")
}
