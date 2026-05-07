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

type mockCreateService struct {
	createFunc func(name string) (*model.Author, error)
}

func (m *mockCreateService) Create(name string) (*model.Author, error) {
	return m.createFunc(name)
}

func TestCreateSuccess(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string) (*model.Author, error) {
			return &model.Author{Id: "1", Name: name}, nil
		},
	}

	req := createRequest{Name: "Test Author"}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/author/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp CreateResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Id != "1" {
		t.Errorf("Expected ID '1', got %s", resp.Id)
	}
	if resp.Status != "Ok" {
		t.Errorf("Expected status 'Ok', got %s", resp.Status)
	}
}

func TestCreateValidationError(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string) (*model.Author, error) {
			return nil, nil
		},
	}

	req := createRequest{Name: ""}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/author/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateDecodeError(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string) (*model.Author, error) {
			return nil, nil
		},
	}

	httpReq := httptest.NewRequest("POST", "/author/create", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateServiceError(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string) (*model.Author, error) {
			return nil, errors.New("author already exists")
		},
	}

	req := createRequest{Name: "Test Author"}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/author/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Create(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
