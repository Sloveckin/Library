package memory

import (
	"errors"
	"sync"

	"Library/internal/model"
	"github.com/google/uuid"
)

var ErrNoSuchAuthor = errors.New("no such author")

type RepositoryInMemory struct {
	rw   sync.RWMutex
	data map[string]*model.Author
}

func NewAuthorRepositoryInMemory() *RepositoryInMemory {
	return &RepositoryInMemory{
		data: make(map[string]*model.Author),
		rw:   sync.RWMutex{},
	}
}

func (s *RepositoryInMemory) Create(name string) (*model.Author, error) {
	author := &model.Author{
		Id:   uuid.NewString(),
		Name: name,
	}

	s.rw.RLock()
	defer s.rw.RUnlock()
	s.data[author.Id] = author

	return author, nil
}

func (s *RepositoryInMemory) Get(id string) (*model.Author, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	author, ok := s.data[id]
	if !ok {
		return nil, ErrNoSuchAuthor
	}

	return author, nil
}

func (s *RepositoryInMemory) Delete(id string) error {
	s.rw.RLock()
	defer s.rw.RUnlock()

	_, ok := s.data[id]
	if !ok {
		return ErrNoSuchAuthor
	}

	delete(s.data, id)

	return nil
}

func (s *RepositoryInMemory) ExistsById(id string) (bool, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	_, ok := s.data[id]

	return ok, nil
}

func (s *RepositoryInMemory) ExistsByName(name string) (bool, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	for _, author := range s.data {
		if author.Name == name {
			return true, nil
		}
	}

	return false, nil
}
