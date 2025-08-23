package cats

import (
	"github.com/go-playground/validator/v10"
)

func NewService(repo *Repository, val *validator.Validate, breeds *Breeds) *Service {
	return &Service{
		Repo:     repo,
		Breeds:   breeds,
		Validate: val,
	}
}
