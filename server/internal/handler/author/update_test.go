package author

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
	updateFunc func(id, name string) (*model.Author, error)
}

func (m *mockUpdateService) Update(id, name string) (*model.Author, error) {
	return m.updateFunc(id, name)
}

func TestUpdateSuccess(t *testing.T) {
	mockService := &mockUpdateService{
		updateFunc: func(id, name string) (*model.Author, error) {
			return &model.Author{Id: id, Name: name}, nil
		},
	}

	req := updateRequest{Id: "1", Name: "Updated Author"}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/author/update", bytes.NewReader(body))
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

func TestUpdateValidationError(t *testing.T) {
	mockService := &mockUpdateService{
		updateFunc: func(id, name string) (*model.Author, error) {
			return nil, nil
		},
	}

	req := updateRequest{Id: "1", Name: ""}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/author/update", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Update(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateServiceError(t *testing.T) {
	mockService := &mockUpdateService{
		updateFunc: func(id, name string) (*model.Author, error) {
			return nil, errors.New("author not found")
		},
	}

	req := updateRequest{Id: "1", Name: "Updated"}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/author/update", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Update(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

var ErrAuthorNotFound error
