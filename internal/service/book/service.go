package service_book

import (
	"Library/internal/model"
	"errors"
)

var (
	BookAlreadyExists = errors.New("book already exists")
	BookNotFound      = errors.New("book not found")
)

type Repository interface {
	Create(name string, authors ...model.Author) (*model.Book, error)
	Get(id string) (*model.Book, error)
	Delete(id string) error
	ExistsById(id string) (bool, error)
	ExistsByName(name string) (bool, error)
}

type ServiceBookImpl struct {
	repository Repository
}

func NewServiceBook(repository Repository) *ServiceBookImpl {
	return &ServiceBookImpl{
		repository: repository,
	}
}

func (s *ServiceBookImpl) Create(name string, authors ...model.Author) (*model.Book, error) {
	exist, err := s.ExistByName(name)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, BookAlreadyExists
	}

	book, err := s.repository.Create(name, authors...)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *ServiceBookImpl) Get(name string) (*model.Book, error) {
	return s.repository.Get(name)
}

func (s *ServiceBookImpl) Delete(id string) error {
	err := s.repository.Delete(id)
	if err != nil {
		return BookNotFound
	}

	return nil
}

func (s *ServiceBookImpl) ExistsById(id string) (bool, error) {
	return s.repository.ExistsById(id)
}

func (s *ServiceBookImpl) ExistByName(name string) (bool, error) {
	return s.repository.ExistsByName(name)
}
