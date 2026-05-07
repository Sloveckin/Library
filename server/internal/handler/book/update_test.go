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

type mockUpdateService struct {
	updateFunc func(id, name string, authors ...model.Author) (*model.Book, error)
}

func (m *mockUpdateService) Update(id, name string, authors ...model.Author) (*model.Book, error) {
	return m.updateFunc(id, name, authors...)
}

func TestUpdateSuccess(t *testing.T) {
	mockService := &mockUpdateService{
		updateFunc: func(id, name string, authors ...model.Author) (*model.Book, error) {
			return &model.Book{Id: id, Name: name, Authors: authors}, nil
		},
	}

	req := updateRequest{
		Id:      "1",
		Name:    "Updated Book",
		Authors: []string{"1", "2"},
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/book/update", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Update(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp updateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Id != "1" {
		t.Errorf("Expected ID '1', got %s", resp.Id)
	}
}

func TestUpdateDecodeError(t *testing.T) {
	mockService := &mockUpdateService{
		updateFunc: func(id, name string, authors ...model.Author) (*model.Book, error) {
			return nil, nil
		},
	}

	httpReq := httptest.NewRequest("PUT", "/book/update", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	handler := Update(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateEmptyId(t *testing.T) {
	mockService := &mockUpdateService{
		updateFunc: func(id, name string, authors ...model.Author) (*model.Book, error) {
			return nil, nil
		},
	}

	req := updateRequest{
		Id:      "",
		Name:    "Updated",
		Authors: []string{"1"},
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/book/update", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Update(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateValidationError(t *testing.T) {
	mockService := &mockUpdateService{
		updateFunc: func(id, name string, authors ...model.Author) (*model.Book, error) {
			return nil, nil
		},
	}

	req := updateRequest{
		Id:      "1",
		Name:    "",
		Authors: []string{},
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/book/update", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Update(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateServiceError(t *testing.T) {
	mockService := &mockUpdateService{
		updateFunc: func(id, name string, authors ...model.Author) (*model.Book, error) {
			return nil, errors.New("book not found")
		},
	}

	req := updateRequest{
		Id:      "1",
		Name:    "Updated",
		Authors: []string{"1"},
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/book/update", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Update(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
