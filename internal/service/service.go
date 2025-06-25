package service

import (
	"github.com/Koyo-os/form-crud-service/internal/entity"
)

type (
	Repository interface{
		Create(*entity.Form) error
		Update(string,string, interface{}) error
		Delete(string) error
		Get(string) (entity.Form, error)
		GetMore(string, interface{}) ([]entity.Form, error)
	}

	Service struct{
		repo Repository
	}
)

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(form *entity.Form) error {
	return s.repo.Create(form)
}

func (s *Service) Update(id string, field string, value interface{}) error {
	return s.repo.Update(id, field, value)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) Get(id string) (entity.Form, error) {
	return s.repo.Get(id)
}

func (s *Service) GetMore(key string, filter interface{}) ([]entity.Form, error) {
	return s.repo.GetMore(key, filter)
}