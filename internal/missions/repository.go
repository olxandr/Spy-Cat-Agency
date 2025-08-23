package missions

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"spy-cat-agency/internal/models"
)

type Repo interface {
	AddTargets(ctx context.Context, missionID int, newTargets []models.Target) ([]models.Target, error)
	AssignCat(ctx context.Context, missionID int, catID int) error
	Create(ctx context.Context, mission *models.Mission) (*models.Mission, error)
	Delete(ctx context.Context, missionID int) error
	DeleteTarget(ctx context.Context, targetID int) error
	Get(ctx context.Context, id int) (*models.Mission, error)
	List(ctx context.Context) (*[]models.Mission, error)
	UpdateAsCompleted(ctx context.Context, missionID int) error
	UpdateTargetNotes(ctx context.Context, targetID int, newNotes string) error
}

var (
	ErrTargetCompleted  = errors.New("Target is completed, unable to edit")
	ErrMIssionCompleted = errors.New("Mission is completed, unable to edit")
	ErrMissionNotFound  = errors.New("Mission not found")
	ErrTooManyTargets   = errors.New("Too many targets")
	ErrCatNotFound      = errors.New("Cat not found")
)

type Repository struct {
	*sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) Create(ctx context.Context, mission *models.Mission) (*models.Mission, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert mission
	insertMissionQuery := `
		INSERT INTO missions (cat_id, is_completed, created_at)
    	VALUES ($1, false, NOW()) RETURNING id
	`

	var missionID int
	err = tx.QueryRow(insertMissionQuery,
		mission.CatID).Scan(&missionID)
	if err != nil {
		slog.Error("Mission insert", "query row", err)
		tx.Rollback()
		return nil, err
	}

	// Insert targets
	inserTargetsQuery := `
		INSERT INTO targets (mission_id, name, country, notes, is_completed)
        VALUES ($1, $2, $3, $4, false)
        RETURNING id
	`
	targets := make([]models.Target, 0, len(mission.Targets))
	stmt, err := tx.Prepare(inserTargetsQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, t := range mission.Targets {
		var targetID int
		if err := stmt.QueryRow(missionID, t.Name, t.Country, t.Notes).Scan(&targetID); err != nil {
			slog.Error("Inserting target", "error", err)
			return nil, err
		}
		targets = append(targets, models.Target{
			ID:          targetID,
			Name:        t.Name,
			Country:     t.Country,
			Notes:       t.Notes,
			IsCompleted: false,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	newMission := &models.Mission{
		ID:          missionID,
		CatID:       new(int),
		IsCompleted: false,
		CreatedAt:   time.Time{},
		Targets:     targets,
	}

	return newMission, nil
}

func (r *Repository) Delete(ctx context.Context, missionID int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if mission is assigned to a cat
	var catID sql.NullInt64
	err = tx.QueryRow(`
        SELECT cat_id FROM missions WHERE id = $1
    `, missionID).Scan(&catID)
	if err != nil {
		return err
	}

	if catID.Valid {
		return errors.New("cannot delete a mission assigned to a cat")
	}

	// Delete the mission (targets will be deleted automatically via ON DELETE CASCADE)
	deleteQuery := `
    	DELETE FROM missions WHERE id = $1
    `
	_, err = tx.Exec(deleteQuery, missionID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAsCompleted(ctx context.Context, missionID int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if mission exists
	existsQuery := `
		SELECT EXISTS(SELECT 1 FROM missions WHERE id = $1)
	`
	var exists bool
	err = tx.QueryRow(existsQuery, missionID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	// Mark the mission as completed
	updateQuery := `
        UPDATE missions SET is_completed = TRUE
        WHERE id = $1
    `
	_, err = tx.Exec(updateQuery, missionID)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateTargetNotes(ctx context.Context, targetID int, newNotes string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if target exists and is not completed
	targetExistsQuery := `
        SELECT is_completed, mission_id FROM targets
        WHERE id = $1
    `
	var isTargetCompleted bool
	var missionID int
	err = tx.QueryRow(targetExistsQuery, targetID).Scan(&isTargetCompleted, &missionID)
	if err != nil {
		return err
	}

	if isTargetCompleted {
		return ErrTargetCompleted
	}

	// Check if parent mission is not completed
	missionNotCompletedQuery := `
        SELECT is_completed FROM missions
        WHERE id = $1
    `
	var isMissionCompleted bool
	err = tx.QueryRow(missionNotCompletedQuery, missionID).Scan(&isMissionCompleted)
	if err != nil {
		return err
	}

	if isMissionCompleted {
		return ErrMIssionCompleted
	}

	//  Update notes
	updateNotesQuery := `
        UPDATE targets SET notes = $1
        WHERE id = $2
    `
	_, err = tx.Exec(updateNotesQuery, newNotes, targetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) DeleteTarget(ctx context.Context, targetID int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if target exists and is not completed
	targetExistsQuery := `
        SELECT is_completed, mission_id FROM targets
        WHERE id = $1
    `
	var isCompleted bool
	var missionID int
	err = tx.QueryRow(targetExistsQuery, targetID).Scan(&isCompleted, &missionID)
	if err != nil {
		return err
	}

	if isCompleted {
		return ErrTargetCompleted
	}

	// Optionally, check if mission exists
	missionExistsQuery := `
		SELECT EXISTS(SELECT 1 FROM missions WHERE id = $1)
	`
	var missionExists bool
	err = tx.QueryRow(missionExistsQuery, missionID).Scan(&missionExists)
	if err != nil {
		return err
	}
	if !missionExists {
		return ErrMissionNotFound
	}

	// Delete the target
	deleteTargetQuery := `
		DELETE FROM targets WHERE id = $1
	`
	_, err = tx.Exec(deleteTargetQuery, targetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) AddTargets(ctx context.Context, missionID int, newTargets []models.Target) ([]models.Target, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Check mission exists and is not completed
	var isMissionCompleted bool
	missionQuery := `
		SELECT is_completed FROM missions WHERE id = $1
	`
	err = tx.QueryRow(missionQuery, missionID).Scan(&isMissionCompleted)
	if err != nil {
		return nil, err
	}

	if isMissionCompleted {
		return nil, ErrMIssionCompleted
	}

	// Count existing targets
	var currentCount int
	countQuery := `
		SELECT COUNT(*) FROM targets WHERE mission_id = $1
	`
	err = tx.QueryRow(countQuery, missionID).Scan(&currentCount)
	if err != nil {
		return nil, err
	}

	if currentCount+len(newTargets) > 3 {
		return nil, ErrTooManyTargets
	}

	// Insert new targets
	insertQuery := `
        INSERT INTO targets (mission_id, name, country, notes, is_completed)
        VALUES ($1, $2, $3, $4, false)
        RETURNING id
    `
	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	insertedTargets := make([]models.Target, 0, len(newTargets))
	for _, t := range newTargets {
		var targetID int
		if err := stmt.QueryRow(missionID, t.Name, t.Country, t.Notes).Scan(&targetID); err != nil {
			return nil, err
		}
		insertedTargets = append(insertedTargets, models.Target{
			ID:          targetID,
			MissionID:   missionID,
			Name:        t.Name,
			Country:     t.Country,
			Notes:       t.Notes,
			IsCompleted: false,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return insertedTargets, nil
}

func (r *Repository) AssignCat(ctx context.Context, missionID int, catID int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check mission exists
	missionExistsQuery := `
		SELECT EXISTS(SELECT 1 FROM missions WHERE id = $1)
	`
	var exists bool
	err = tx.QueryRow(missionExistsQuery, missionID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrMissionNotFound
	}

	// Check cat exists
	catExistsQuery := `
		SELECT EXISTS(SELECT 1 FROM cats WHERE id = $1)
	`
	err = tx.QueryRow(catExistsQuery, catID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCatNotFound
	}

	// Assign cat to mission
	updateQuery := `
        UPDATE missions SET cat_id = $1
        WHERE id = $2
    `
	_, err = tx.Exec(updateQuery, catID, missionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) List(ctx context.Context) (*[]models.Mission, error) {
	query := `
		SELECT id, cat_id, is_completed, created_at FROM missions
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		if err := rows.Scan(&mission.ID, &mission.CatID, &mission.IsCompleted, &mission.CreatedAt); err != nil {
			return nil, err
		}
		missions = append(missions, mission)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &missions, nil
}

func (r *Repository) Get(ctx context.Context, id int) (*models.Mission, error) {
	query := `
		SELECT id, cat_id, is_completed, created_at FROM missions
		WHERE id = $1
	`

	var mission models.Mission
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&mission.ID, &mission.CatID, &mission.IsCompleted, &mission.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &mission, nil
}
