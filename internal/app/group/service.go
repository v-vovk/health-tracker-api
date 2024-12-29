package group

type Service interface {
	GetAll(limit, offset int) ([]Group, int64, error)
	GetByID(id string) (*Group, error)
	Create(group *Group) error
	Update(group *Group) error
	Delete(id string) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) GetAll(limit, offset int) ([]Group, int64, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *serviceImpl) GetByID(id string) (*Group, error) {
	group, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *serviceImpl) Create(group *Group) error {
	return s.repo.Create(group)
}

func (s *serviceImpl) Update(group *Group) error {
	existingGroup, err := s.repo.GetByID(group.ID)
	if err != nil {
		return err
	}

	group.CreatedAt = existingGroup.CreatedAt

	return s.repo.Update(group)
}

func (s *serviceImpl) Delete(id string) error {
	return s.repo.Delete(id)
}
