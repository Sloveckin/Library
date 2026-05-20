package book

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/model"
)

const expectedGetBookStatusMsg = "Expected status %d, got %d"

type mockGetService struct {
	getFunc func(id string) (*model.Book, error)
}

func (m *mockGetService) Get(id string) (*model.Book, error) {
	return m.getFunc(id)
}

func TestGetSuccess(t *testing.T) {
	mockService := &mockGetService{
		getFunc: func(id string) (*model.Book, error) {
			return &model.Book{
				Id:   id,
				Name: "Test Book",
				Authors: []model.Author{
					{Id: "1", Name: "Author 1"},
				},
			}, nil
		},
	}

	httpReq := httptest.NewRequest("GET", "/book/get?id=1", nil)
	w := httptest.NewRecorder()

	handler := Get(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf(expectedGetBookStatusMsg, http.StatusCreated, w.Code)
	}
}

func TestGetMissingId(t *testing.T) {
	mockService := &mockGetService{
		getFunc: func(id string) (*model.Book, error) {
			return nil, nil
		},
	}

	httpReq := httptest.NewRequest("GET", "/book/get", nil)
	w := httptest.NewRecorder()

	handler := Get(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf(expectedGetBookStatusMsg, http.StatusBadRequest, w.Code)
	}
}

func TestGetServiceError(t *testing.T) {
	mockService := &mockGetService{
		getFunc: func(id string) (*model.Book, error) {
			return nil, errors.New("not found")
		},
	}

	httpReq := httptest.NewRequest("GET", "/book/get?id=1", nil)
	w := httptest.NewRecorder()

	handler := Get(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf(expectedGetBookStatusMsg, http.StatusNotFound, w.Code)
}
}
