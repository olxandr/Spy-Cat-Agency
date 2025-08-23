package cats

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"spy-cat-agency/internal/models"

	"github.com/go-playground/validator/v10"
)

type Breed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Breeds struct {
	Api   string
	Cache map[string]struct{}
	mu    sync.RWMutex
}

type Service struct {
	Repo
	*Breeds
	Validate *validator.Validate
}

func NewBreeds(breedsApi string) (*Breeds, error) {
	b := &Breeds{
		Api:   breedsApi,
		Cache: map[string]struct{}{},
		mu:    sync.RWMutex{},
	}
	if err := b.Fetch(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Breeds) Exists(fl validator.FieldLevel) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	breedName := fl.Field().String()
	_, found := b.Cache[breedName]
	return found
}

func (b *Breeds) Fetch() error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(b.Api)
	if err != nil {
		return fmt.Errorf("Error executing breeds API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading breeds API response body: %w", err)
	}

	var breeds []Breed
	if err := json.Unmarshal(body, &breeds); err != nil {
		return fmt.Errorf("Error unmarshaling breeds API response: %w", err)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.Cache = make(map[string]struct{})
	for _, breed := range breeds {
		b.Cache[breed.Name] = struct{}{}
	}

	slog.Info("Breeds cache populated", "breeds", len(b.Cache))
	return nil
}

func (s *Service) ValidateCat(cat *models.Cat) (validator.ValidationErrors, bool) {
	return cat.Valid(s.Validate)
}

func (s *Service) Create(ctx context.Context, cat *models.Cat) (int64, error) {
	id, err := s.Repo.Create(ctx, cat)
	if err != nil {
		return 0, fmt.Errorf("Failed to insert a cat: %w", err)
	}

	slog.Info("New spy cat is created", "id", id, "name", cat.Name)

	return id, nil
}

func (s *Service) Remove(ctx context.Context, id int) error {
	return s.Repo.Remove(ctx, id)
}

func (s *Service) UpdateSalary(ctx context.Context, cat *models.Cat) (*models.Cat, error) {
	return s.Repo.UpdateSalary(ctx, cat)
}

func (s *Service) List(ctx context.Context) ([]models.Cat, error) {
	return s.Repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int) (*models.Cat, error) {
	return s.Repo.Get(ctx, id)
}
