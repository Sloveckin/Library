package service_book

import (
	"errors"

	"Library/internal/model"
)

var (
	ErrBookAlreadyExists = errors.New("book already exists")
	ErrBookNotFound      = errors.New("book not found")
	ErrAuthorNotFound    = errors.New("author not found")
)

//go:generate mockery --name=Repository
type Repository interface {
	Create(name string, authors ...model.Author) (*model.Book, error)
	Get(id string) (*model.Book, error)
	Delete(id string) error
	ExistsById(id string) (bool, error)
	ExistsByName(name string) (bool, error)
}

//go:generate mockery --name=AuthorService
type AuthorService interface {
	ExistsById(id string) (bool, error)
}

type ServiceBookImpl struct {
	repository    Repository
	authorService AuthorService
}

func NewServiceBook(repository Repository, service AuthorService) *ServiceBookImpl {
	return &ServiceBookImpl{
		repository:    repository,
		authorService: service,
	}
}

func (s *ServiceBookImpl) Create(name string, authors ...model.Author) (*model.Book, error) {
	exist, err := s.ExistByName(name)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, ErrBookAlreadyExists
	}

	for _, author := range authors {
		ok, err := s.authorService.ExistsById(author.Id)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, ErrAuthorNotFound
		}
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
		return ErrBookNotFound
	}

	return nil
}

func (s *ServiceBookImpl) ExistsById(id string) (bool, error) {
	return s.repository.ExistsById(id)
}

func (s *ServiceBookImpl) ExistByName(name string) (bool, error) {
	return s.repository.ExistsByName(name)
}
