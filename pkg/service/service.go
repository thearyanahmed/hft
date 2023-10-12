package service

import "github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"

type Service struct {
	configRepo ConfigRepository
}

type ConfigRepository interface {
	Store(entity schema.ConfigMap) error
	Find(options *schema.FilterOptions) ([]schema.ConfigMap, error)
	Update(entity schema.ConfigMap) error
	Delete(entity schema.ConfigMap) error
}

func NewService(configRepo ConfigRepository) *Service {
	return &Service{
		configRepo: configRepo,
	}
}

func NewInMemoryRepository() *Service {
	return &Service{}
}

func (s *Service) Store(entity schema.ConfigMap) error {
	return s.configRepo.Store(entity)
}

func (s *Service) Find(options *schema.FilterOptions) ([]schema.ConfigMap, error) {
	return s.configRepo.Find(options)
}

func (s *Service) Update(entity schema.ConfigMap) error {
	return s.configRepo.Update(entity)
}

func (s *Service) Delete(entity schema.ConfigMap) error {
	return nil
}
