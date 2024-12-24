package food

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) GetAllFoods(limit, offset int) ([]Food, int64, error) {
	return s.Repo.GetAllFoods(limit, offset)
}

func (s *Service) GetFoodByID(id string) (*Food, error) {
	return s.Repo.GetFoodByID(id)
}

func (s *Service) CreateFood(food *Food) error {
	return s.Repo.CreateFood(food)
}

func (s *Service) UpdateFood(food *Food) error {
	return s.Repo.UpdateFood(food)
}

func (s *Service) DeleteFood(id string) error {
	return s.Repo.DeleteFood(id)
}
