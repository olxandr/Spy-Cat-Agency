package models

// type Target struct {
// 	ID          int    `json:"id" validate:"required,gt=0"`
// 	MissionID   int    `json:"mission_id"`
// 	Name        string `json:"name" validate:"required,min=2,max=50"`
// 	Country     string `json:"country" validate:"required,min=2,max=50"`
// 	Notes       string `json:"notes"`
// 	IsCompleted bool   `json:"is_completed" validate:"required"`
// }

type Target struct {
	ID          int    `json:"id,omitempty"`
	MissionID   int    `json:"mission_id,omitempty"`
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Country     string `json:"country" validate:"required,min=2,max=50"`
	Notes       string `json:"notes,omitempty"`
	IsCompleted bool   `json:"is_completed,omitempty"`
}
