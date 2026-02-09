package memory

import (
	"Library/internal/model"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	NoSuchBook = errors.New("no such book")
)

type RepositoryImpl struct {
	rw   sync.RWMutex
	data map[string]*model.Book
}

func NewRepositoryInMemory() *RepositoryImpl {
	return &RepositoryImpl{
		rw:   sync.RWMutex{},
		data: make(map[string]*model.Book),
	}
}

func (repo *RepositoryImpl) Create(name string, author ...model.Author) (*model.Book, error) {

	book := &model.Book{
		Id:      uuid.NewString(),
		Name:    name,
		Authors: author,
	}

	repo.rw.Lock()
	repo.data[book.Id] = book
	repo.rw.Unlock()

	return book, nil
}

func (repo *RepositoryImpl) Get(id string) (*model.Book, error) {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	if book, ok := repo.data[id]; ok {
		return book, nil
	}

	return nil, NoSuchBook
}

func (repo *RepositoryImpl) ExistsById(id string) (bool, error) {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	_, ok := repo.data[id]

	return ok, nil
}

func (repo *RepositoryImpl) ExistsByName(name string) (bool, error) {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	for _, book := range repo.data {
		if book.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func (repo *RepositoryImpl) Delete(id string) error {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	if _, ok := repo.data[id]; ok {
		delete(repo.data, id)
		return nil
	}

	return NoSuchBook
}
