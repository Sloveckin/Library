package service_author

import (
	"errors"

	"Library/internal/model"
)

var (
	ErrAuthorExists    = errors.New("author already exists")
	ErrAuthorNotExists = errors.New("author not exists")
)

//go:generate mockery --name=AuthorRepository
type AuthorRepository interface {
	Create(name string) (*model.Author, error)
	Get(id string) (*model.Author, error)
	GetByName(name string) (*model.Author, error)
	Update(id, name string) (*model.Author, error)
	Delete(id string) error
	ExistsById(id string) (bool, error)
	ExistsByName(name string) (bool, error)
}

type AuthorServiceImpl struct {
	authorRepository AuthorRepository
}

func NewAuthorServiceImpl(repo AuthorRepository) *AuthorServiceImpl {
	return &AuthorServiceImpl{
		authorRepository: repo,
	}
}

func (s *AuthorServiceImpl) Create(name string) (*model.Author, error) {
	exist, err := s.ExistsByName(name)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, ErrAuthorExists
	}

	author, err := s.authorRepository.Create(name)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorServiceImpl) Get(id string) (*model.Author, error) {
	return s.authorRepository.Get(id)
}

func (s *AuthorServiceImpl) Delete(id string) error {
	exists, err := s.ExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return ErrAuthorNotExists
	}

	err = s.authorRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthorServiceImpl) Update(id, name string) (*model.Author, error) {
	exists, err := s.ExistsById(id)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrAuthorNotExists
	}

	// If name already exists and belongs to a different author -> error
	existByName, err := s.ExistsByName(name)
	if err != nil {
		return nil, err
	}

	if existByName {
		// fetch author with that name to check id
		authorWithName, err := s.authorRepository.GetByName(name)
		if err != nil {
			return nil, err
		}

		if authorWithName != nil && authorWithName.Id != id {
			return nil, ErrAuthorExists
		}
	}

	updated, err := s.authorRepository.Update(id, name)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *AuthorServiceImpl) ExistsByName(name string) (bool, error) {
	return s.authorRepository.ExistsByName(name)
}

func (s *AuthorServiceImpl) ExistsById(id string) (bool, error) {
	return s.authorRepository.ExistsById(id)
}
