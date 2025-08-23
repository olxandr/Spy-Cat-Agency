package models

import (
	"time"
)

type Mission struct {
	ID          int       `json:"id,omitempty"`
	CatID       int       `json:"cat_id,omitempty"`
	IsCompleted bool      `json:"is_completed,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitzero"`
	Targets     []Target  `json:"targets,omitempty"`
}
