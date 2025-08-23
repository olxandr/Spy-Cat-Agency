package models

type Cat struct {
	ID                int64   `json:"id,omitempty"`
	Name              string  `json:"name,omitempty"`
	YearsOfExperience int8    `json:"yoe,omitempty"`
	Breed             string  `json:"breed,omitempty"`
	Salary            float64 `json:"salary,omitempty"`
	MissionID         int     `json:"mission_id,omitempty"`
}
