package book

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/model"
)

const (
	getAllBooksPath    = "/book/all"
	expectedAllBookStatusMsg  = "Expected status %d, got %d"
)

type mockAllService struct {
	getAllFunc func() ([]model.Book, error)
}

func (m *mockAllService) GetAll() ([]model.Book, error) {
	return m.getAllFunc()
}

func TestAllSuccess(t *testing.T) {
	mockService := &mockAllService{
		getAllFunc: func() ([]model.Book, error) {
			return []model.Book{
				{
					Id:   "1",
					Name: "Book 1",
					Authors: []model.Author{
						{Id: "1", Name: "Author 1"},
					},
				},
				{
					Id:   "2",
					Name: "Book 2",
					Authors: []model.Author{
						{Id: "2", Name: "Author 2"},
					},
				},
			}, nil
		},
	}

	httpReq := httptest.NewRequest("GET", getAllBooksPath, nil)
	w := httptest.NewRecorder()

	handler := All(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf(expectedAllBookStatusMsg, http.StatusOK, w.Code)
	}
}

func TestAllEmptyList(t *testing.T) {
	mockService := &mockAllService{
		getAllFunc: func() ([]model.Book, error) {
			return []model.Book{}, nil
		},
	}

	httpReq := httptest.NewRequest("GET", getAllBooksPath, nil)
	w := httptest.NewRecorder()

	handler := All(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf(expectedAllBookStatusMsg, http.StatusOK, w.Code)
	}
}

func TestAllServiceError(t *testing.T) {
	mockService := &mockAllService{
		getAllFunc: func() ([]model.Book, error) {
			return nil, errors.New("database error")
		},
	}

	httpReq := httptest.NewRequest("GET", getAllBooksPath, nil)
	w := httptest.NewRecorder()

	handler := All(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusInternalServerError {
		t.Errorf(expectedAllBookStatusMsg, http.StatusInternalServerError, w.Code)
}
}
