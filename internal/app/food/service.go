package food

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetAll(limit, offset int) ([]Food, int64, error) {
	return s.Repo.GetAll(limit, offset)
}

func (s *Service) GetByID(id string) (*Food, error) {
	return s.Repo.GetByID(id)
}

func (s *Service) Create(food *Food) error {
	return s.Repo.Create(food)
}

func (s *Service) Update(food *Food) error {
	existingFood, err := s.Repo.GetByID(food.ID)
	if err != nil {
		return err
	}

	food.CreatedAt = existingFood.CreatedAt

	return s.Repo.Update(food)
}

func (s *Service) Delete(id string) error {
	return s.Repo.Delete(id)
}
