package models

type Target struct {
	ID          int    `json:"id,omitempty"`
	MissionID   int    `json:"mission_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Country     string `json:"country,omitempty"`
	Notes       string `json:"notes,omitempty"`
	IsCompleted bool   `json:"is_completed,omitempty"`
}
