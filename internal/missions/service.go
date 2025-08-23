package missions

func NewService(repo *Repository) *Service {
	return &Service{
		Repo: repo,
	}
}
