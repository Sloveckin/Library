package service_book

import (
	"errors"
	"testing"

	"server/internal/model"
)

type mockRepository struct {
	createFunc      func(name string, authors ...model.Author) (*model.Book, error)
	getFunc         func(id string) (*model.Book, error)
	getByNameFunc   func(name string) (*model.Book, error)
	updateFunc      func(id, name string, authors ...model.Author) (*model.Book, error)
	deleteFunc      func(id string) error
	existsByIdFunc  func(id string) (bool, error)
	existsByNameFunc func(name string) (bool, error)
}

func (m *mockRepository) Create(name string, authors ...model.Author) (*model.Book, error) {
	return m.createFunc(name, authors...)
}

func (m *mockRepository) Get(id string) (*model.Book, error) {
	return m.getFunc(id)
}

func (m *mockRepository) GetByName(name string) (*model.Book, error) {
	return m.getByNameFunc(name)
}

func (m *mockRepository) Update(id, name string, authors ...model.Author) (*model.Book, error) {
	return m.updateFunc(id, name, authors...)
}

func (m *mockRepository) Delete(id string) error {
	return m.deleteFunc(id)
}

func (m *mockRepository) ExistsById(id string) (bool, error) {
	return m.existsByIdFunc(id)
}

func (m *mockRepository) ExistsByName(name string) (bool, error) {
	return m.existsByNameFunc(name)
}

type mockAuthorService struct {
	existsByIdFunc func(id string) (bool, error)
}

func (m *mockAuthorService) ExistsById(id string) (bool, error) {
	return m.existsByIdFunc(id)
}

func TestCreateSuccess(t *testing.T) {
	mockRepo := &mockRepository{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return &model.Book{Id: "1", Name: name, Authors: authors}, nil
		},
		existsByNameFunc: func(name string) (bool, error) {
			return false, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
	}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	book, err := service.Create("Test Book", model.Author{Id: "1"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if book.Id != "1" {
		t.Errorf("Expected ID '1', got %s", book.Id)
	}
}

func TestCreateBookExists(t *testing.T) {
	mockRepo := &mockRepository{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return nil, nil
		},
		existsByNameFunc: func(name string) (bool, error) {
			return true, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
	}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Create("Test Book")

	if err != ErrBookAlreadyExists {
		t.Errorf("Expected ErrBookAlreadyExists, got %v", err)
	}
}

func TestCreateAuthorNotFound(t *testing.T) {
	mockRepo := &mockRepository{
		existsByNameFunc: func(name string) (bool, error) {
			return false, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{
		existsByIdFunc: func(id string) (bool, error) {
			return false, nil
		},
	}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Create("Test Book", model.Author{Id: "999"})

	if err != ErrAuthorNotFound {
		t.Errorf("Expected ErrAuthorNotFound, got %v", err)
	}
}

func TestGet(t *testing.T) {
	mockRepo := &mockRepository{
		getFunc: func(id string) (*model.Book, error) {
			return &model.Book{Id: id, Name: "Test"}, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	book, err := service.Get("1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if book.Id != "1" {
		t.Errorf("Expected ID '1', got %s", book.Id)
	}
}

func TestDelete(t *testing.T) {
	mockRepo := &mockRepository{
		deleteFunc: func(id string) error {
			return nil
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	err := service.Delete("1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteError(t *testing.T) {
	mockRepo := &mockRepository{
		deleteFunc: func(id string) error {
			return errors.New("delete failed")
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	err := service.Delete("1")

	if err != ErrBookNotFound {
		t.Errorf("Expected ErrBookNotFound, got %v", err)
	}
}

func TestUpdateSuccess(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
		existsByNameFunc: func(name string) (bool, error) {
			return false, nil
		},
		updateFunc: func(id, name string, authors ...model.Author) (*model.Book, error) {
			return &model.Book{Id: id, Name: name, Authors: authors}, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
	}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	book, err := service.Update("1", "Updated", model.Author{Id: "1"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if book.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got %s", book.Name)
	}
}

func TestUpdateNotFound(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return false, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Update("1", "Updated")

	if err != ErrBookNotFound {
		t.Errorf("Expected ErrBookNotFound, got %v", err)
	}
}

func TestExistsById(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	exists, err := service.ExistsById("1")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("Expected exists=true, got false")
	}
}

func TestExistByName(t *testing.T) {
	mockRepo := &mockRepository{
		existsByNameFunc: func(name string) (bool, error) {
			return true, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	exists, err := service.ExistByName("Test")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("Expected exists=true, got false")
	}
}

func TestCreateExistsByNameError(t *testing.T) {
	mockRepo := &mockRepository{
		existsByNameFunc: func(name string) (bool, error) {
			return false, errors.New("db error")
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Create("Test Book")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreateAuthorServiceError(t *testing.T) {
	mockRepo := &mockRepository{
		existsByNameFunc: func(name string) (bool, error) {
			return false, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{
		existsByIdFunc: func(id string) (bool, error) {
			return false, errors.New("service error")
		},
	}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Create("Test Book", model.Author{Id: "1"})

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreateRepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return nil, errors.New("create failed")
		},
		existsByNameFunc: func(name string) (bool, error) {
			return false, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
	}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Create("Test Book", model.Author{Id: "1"})

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestExistsByIdError(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return false, errors.New("db error")
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.ExistsById("1")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestExistByNameError(t *testing.T) {
	mockRepo := &mockRepository{
		existsByNameFunc: func(name string) (bool, error) {
			return false, errors.New("db error")
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.ExistByName("Test")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUpdateExistsByIdError(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return false, errors.New("db error")
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Update("1", "Updated")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUpdateExistsByNameError(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
		existsByNameFunc: func(name string) (bool, error) {
			return false, errors.New("db error")
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Update("1", "Updated")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUpdateBookNameExists(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
		existsByNameFunc: func(name string) (bool, error) {
			return true, nil
		},
		getByNameFunc: func(name string) (*model.Book, error) {
			return &model.Book{Id: "999", Name: name}, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Update("1", "Updated")

	if err != ErrBookAlreadyExists {
		t.Errorf("Expected ErrBookAlreadyExists, got %v", err)
	}
}

func TestUpdateGetByNameError(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
		existsByNameFunc: func(name string) (bool, error) {
			return true, nil
		},
		getByNameFunc: func(name string) (*model.Book, error) {
			return nil, errors.New("get failed")
		},
	}

	mockAuthorSvc := &mockAuthorService{}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Update("1", "Updated")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestUpdateAuthorNotFound(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
		existsByNameFunc: func(name string) (bool, error) {
			return false, nil
		},
	}

	mockAuthorSvc := &mockAuthorService{
		existsByIdFunc: func(id string) (bool, error) {
			return false, nil
		},
	}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Update("1", "Updated", model.Author{Id: "999"})

	if err != ErrAuthorNotFound {
		t.Errorf("Expected ErrAuthorNotFound, got %v", err)
	}
}

func TestUpdateRepositoryError(t *testing.T) {
	mockRepo := &mockRepository{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
		existsByNameFunc: func(name string) (bool, error) {
			return false, nil
		},
		updateFunc: func(id, name string, authors ...model.Author) (*model.Book, error) {
			return nil, errors.New("update failed")
		},
	}

	mockAuthorSvc := &mockAuthorService{
		existsByIdFunc: func(id string) (bool, error) {
			return true, nil
		},
	}

	service := NewServiceBook(mockRepo, mockAuthorSvc)
	_, err := service.Update("1", "Updated", model.Author{Id: "1"})

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
