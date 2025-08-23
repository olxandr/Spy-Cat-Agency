package missions

import "github.com/go-playground/validator/v10"

func NewService(repo *Repository, valid *validator.Validate) *Service {
	return &Service{
		Repo:     repo,
		Validate: valid,
	}
}
