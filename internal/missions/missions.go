package missions

import (
	"context"
	"errors"

	"spy-cat-agency/internal/models"
)

var ErrInvalidID = errors.New("ID can't be 0")

type Service struct {
	Repo
}

func (s *Service) Create(ctx context.Context, mission *models.Mission) (*models.Mission, error) {
	return s.Repo.Create(ctx, mission)
}

func (s *Service) Delete(ctx context.Context, missionID int) error {
	return s.Repo.Delete(ctx, missionID)
}

func (s *Service) UpdateAsCompleted(ctx context.Context, missionID int) error {
	return s.Repo.UpdateAsCompleted(ctx, missionID)
}

func (s *Service) UpdateTargetNotes(ctx context.Context, targetID int, notes string) error {
	return s.Repo.UpdateTargetNotes(ctx, targetID, notes)
}

func (s *Service) DeleteTarget(ctx context.Context, targetID int) error {
	return s.Repo.DeleteTarget(ctx, targetID)
}

func (s *Service) AddTargets(ctx context.Context, missionID int, targets []models.Target) ([]models.Target, error) {
	return s.Repo.AddTargets(ctx, missionID, targets)
}

func (s *Service) AssignCat(ctx context.Context, missionID int, catID int) error {
	return s.Repo.AssignCat(ctx, missionID, catID)
}

func (s *Service) List(ctx context.Context) (*[]models.Mission, error) {
	return s.Repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int) (*models.Mission, error) {
	return s.Repo.Get(ctx, id)
}
