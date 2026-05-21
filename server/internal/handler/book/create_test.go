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

const (
	createBookPath              = "/book/create"
	expectedCreateBookStatusMsg = "Expected status %d, got %d"
)

type mockCreateService struct {
	createFunc func(name string, authors ...model.Author) (*model.Book, error)
}

func (m *mockCreateService) Create(name string, authors ...model.Author) (*model.Book, error) {
	return m.createFunc(name, authors...)
}

func runCreateRequest(t *testing.T, service createService, body []byte) *httptest.ResponseRecorder {
	t.Helper()

	httpReq := httptest.NewRequest(http.MethodPost, createBookPath, bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler := Create(service)
	handler(w, httpReq)

	return w
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
	w := runCreateRequest(t, mockService, body)

	if w.Code != http.StatusCreated {
		t.Errorf(expectedCreateBookStatusMsg, http.StatusCreated, w.Code)
	}

	var resp createResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
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
	w := runCreateRequest(t, mockService, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf(expectedCreateBookStatusMsg, http.StatusBadRequest, w.Code)
	}
}

func TestCreateDecodeError(t *testing.T) {
	mockService := &mockCreateService{
		createFunc: func(name string, authors ...model.Author) (*model.Book, error) {
			return nil, nil
		},
	}

	w := runCreateRequest(t, mockService, []byte("invalid"))

	if w.Code != http.StatusBadRequest {
		t.Errorf(expectedCreateBookStatusMsg, http.StatusBadRequest, w.Code)
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
	w := runCreateRequest(t, mockService, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf(expectedCreateBookStatusMsg, http.StatusBadRequest, w.Code)
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
	w := runCreateRequest(t, mockService, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf(expectedCreateBookStatusMsg, http.StatusBadRequest, w.Code)
	}
}
