package service_book

import (
	"Library/internal/model"
	"Library/internal/repo/InMemory"
	"errors"
)

var (
	BookAlreadyExists = errors.New("book already exists")
	BookNotFound      = errors.New("book not found")
)

type Repository interface {
	Create(name string, authors ...model.Author) (*model.Book, error)
	Get(name string) (*model.Book, error)
	Delete(name string) error
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

func (s *ServiceBookImpl) Delete(name string) error {
	exists, err := s.ExistByName(name)
	if err != nil {
		return err
	}

	if !exists {
		return BookNotFound
	}

	return s.repository.Delete(name)
}

func (s *ServiceBookImpl) ExistByName(name string) (bool, error) {
	_, err := s.repository.Get(name)
	if err != nil {
		if errors.Is(err, InMemory.NoSuchBook) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
