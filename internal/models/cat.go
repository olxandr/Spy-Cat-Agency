package models

import (
	"github.com/go-playground/validator/v10"
)

type Cat struct {
	ID                int64   `json:"id,omitempty"`
	Name              string  `json:"name" validate:"required,min=2,max=50"`
	YearsOfExperience int8    `json:"yoe" validate:"lt=100"`
	Breed             string  `json:"breed" validate:"required,validBreed,min=2,max=50"`
	Salary            float64 `json:"salary" validate:"required"`
	MissionID         *int    `json:"mission_id,omitempty"`
}

func (c *Cat) Valid(v *validator.Validate) (validator.ValidationErrors, bool) {
	if err := v.Struct(c); err != nil {
		validationErrs, ok := err.(validator.ValidationErrors)
		if ok {
			return validationErrs, false
		}
	}

	return nil, true
}
