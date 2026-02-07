package InMemory

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
	repo.data[book.Name] = book
	repo.rw.Unlock()

	return book, nil
}

func (repo *RepositoryImpl) Get(name string) (*model.Book, error) {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	if book, ok := repo.data[name]; ok {
		return book, nil
	}

	return nil, NoSuchBook
}

func (repo *RepositoryImpl) Delete(name string) error {
	repo.rw.Lock()
	defer repo.rw.Unlock()

	if _, ok := repo.data[name]; ok {
		delete(repo.data, name)
		return nil
	}

	return NoSuchBook
}
