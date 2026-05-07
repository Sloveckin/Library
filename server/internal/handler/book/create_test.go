package book

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/model"
)

type mockCreateService struct {
	createFunc func(name string, authors ...model.Author) (*model.Book, error)
}

func (m *mockCreateService) Create(name string, authors ...model.Author) (*model.Book, error) {
	return m.createFunc(name, authors...)
}

func TestCreateSuccess(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return &model.Book{Id: "1", Name: name, Authors: authors}, nil
		},
	}

	req := createRequest{
		Name:     "Test Book",
		AuthorId: []string{"1", "2"},
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/book/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp createResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Id != "1" {
		t.Errorf("Expected ID '1', got %s", resp.Id)
	}
}

func TestCreateValidationError(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return nil, nil
		},
	}

	req := createRequest{Name: "", AuthorId: []string{}}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/book/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateDecodeError(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return nil, nil
		},
	}

	httpReq := httptest.NewRequest("POST", "/book/create", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateServiceError(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return nil, errors.New("book exists")
		},
	}

	req := createRequest{
		Name:     "Test Book",
		AuthorId: []string{"1"},
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/book/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateEmptyName(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return nil, errors.New("book exists")
		},
	}

	req := createRequest{
		Name:     "   ",
		AuthorId: []string{"1"},
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/book/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
