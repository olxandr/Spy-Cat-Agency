package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Mission struct {
	ID          int       `json:"id,omitempty"`
	CatID       *int      `json:"cat_id,omitempty"`
	IsCompleted bool      `json:"is_completed,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitzero"`
	Targets     []Target  `json:"targets,omitempty"`
}

func (m *Mission) Valid(v *validator.Validate) (validator.ValidationErrors, bool) {
	if err := v.Struct(m); err != nil {
		validationErrs, ok := err.(validator.ValidationErrors)
		if ok {
			return validationErrs, false
		}
	}

	return nil, true
}
