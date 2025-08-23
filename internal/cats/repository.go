package cats

import (
	"context"
	"database/sql"
	"log/slog"

	"spy-cat-agency/internal/models"
)

type Repo interface {
	Create(ctx context.Context, cat *models.Cat) (int64, error)
	Get(ctx context.Context, id int) (*models.Cat, error)
	List(ctx context.Context) ([]models.Cat, error)
	Remove(ctx context.Context, id int) error
	UpdateSalary(ctx context.Context, cat *models.Cat) (*models.Cat, error)
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

type Repository struct {
	DB *sql.DB
}

func (s *Repository) Create(ctx context.Context, cat *models.Cat) (int64, error) {
	query := `
		INSERT INTO cats (name, years_of_experience, breed, salary)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int64
	if err := s.DB.QueryRow(query,
		cat.Name,
		cat.YearsOfExperience,
		cat.Breed,
		cat.Salary,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Repository) Remove(ctx context.Context, id int) error {
	query := `
		DELETE FROM cats WHERE id = $1
	`
	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		slog.Info("Remove cat", "exec context", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Repository) UpdateSalary(ctx context.Context, cat *models.Cat) (*models.Cat, error) {
	updateQuery := `
		UPDATE cats SET salary = $1 WHERE id = $2
	`

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.Exec(updateQuery, cat.Salary, cat.ID)
	if err != nil {
		slog.Error("UpdateSalary", "update query exec error", err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	getQuery := `
		SELECT salary FROM cats WHERE id = $1
	`

	updatedCat := *cat
	if err = tx.QueryRow(getQuery, cat.ID).Scan(&updatedCat.Salary); err != nil {
		slog.Error("Get updated salary cat", "error", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("Update salary", "commit transaction", err)
		return nil, err
	}

	return &updatedCat, nil
}

func (s *Repository) List(ctx context.Context) ([]models.Cat, error) {
	query := `
		SELECT id, name, years_of_experience, breed, salary
		FROM cats
	`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		slog.Error("List cats", "query context", err)
		return nil, err
	}
	defer rows.Close()

	var cats []models.Cat
	for rows.Next() {
		var cat models.Cat
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary); err != nil {
			slog.Error("List cats", "rows scan", err)
			return nil, err
		}
		cats = append(cats, cat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cats, nil
}

func (s *Repository) Get(ctx context.Context, id int) (*models.Cat, error) {
	query := `
		SELECT id, name, years_of_experience, breed, salary
		FROM cats
		WHERE id = $1
	`
	var cat models.Cat
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary)
	if err != nil {
		slog.Error("Get cat", "query exec", err)
		return nil, err
	}

	return &cat, nil
}
