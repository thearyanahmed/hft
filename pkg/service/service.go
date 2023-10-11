package service

import "github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/schema"

type Service struct {
	configRepo ConfigRepository
}

type ConfigRepository interface {
	Store(entity schema.ConfigMap) error
	Find() ([]schema.ConfigMap, error)
	Update() error
	Delete() error
}

func NewService(configRepo ConfigRepository) *Service {
	return &Service{
		configRepo: configRepo,
	}
}

func NewInMemoryRepository() *Service {
	return &Service{}
}

func (s *Service) Store() error {
	return nil
}

// @todo note use filters and limit
func (s *Service) Find() ([]schema.ConfigMap, error) {
	return s.configRepo.Find()
}
func (s *Service) Update() error {
	return nil
}
func (s *Service) Delete() error {
	return nil
}
