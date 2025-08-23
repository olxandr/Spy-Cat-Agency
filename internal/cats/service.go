package cats

func NewService(repo *Repository, breeds *Breeds) *Service {
	return &Service{
		Repo:   repo,
		Breeds: breeds,
	}
}
