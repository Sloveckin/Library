package service_author

import (
	"Library/internal/model"
	"errors"
)

var (
	AuthorExists    = errors.New("author already exists")
	AuthorNotExists = errors.New("author not exists")
)

type AuthorRepository interface {
	Create(name string) (*model.Author, error)
	Get(id string) (string, error)
	Delete(id string) error
	ExistsById(id string) (bool, error)
	ExistsByName(name string) (bool, error)
}

type AuthorServiceImpl struct {
	authorRepository AuthorRepository
}

func NewAuthorServiceImpl() *AuthorServiceImpl {
	return &AuthorServiceImpl{}
}

func (s *AuthorServiceImpl) Create(name string) (*model.Author, error) {
	exist, err := s.ExistsByName(name)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, AuthorExists
	}

	author, err := s.authorRepository.Create(name)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorServiceImpl) Get(id string) (string, error) {
	return s.authorRepository.Get(id)
}

func (s *AuthorServiceImpl) Delete(id string) error {
	err := s.authorRepository.Delete(id)
	if err != nil {
		return AuthorNotExists
	}

	return nil
}

func (s *AuthorServiceImpl) ExistsByName(name string) (bool, error) {
	return s.authorRepository.ExistsByName(name)
}

func (s *AuthorServiceImpl) ExistsById(id string) (bool, error) {
	return s.authorRepository.ExistsById(id)
}
