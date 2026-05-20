package author

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/model"
)

const (
	getAuthorPath      = "/author/get?id=1"
	expectedGetAuthorStatusMsg  = "Expected status %d, got %d"
)

type mockGetService struct {
	getFunc func(id string) (*model.Author, error)
}

func (m *mockGetService) Get(id string) (*model.Author, error) {
	return m.getFunc(id)
}

func TestGetSuccess(t *testing.T) {
	mockService := &mockGetService{
		getFunc: func(id string) (*model.Author, error) {
			return &model.Author{Id: id, Name: "Test"}, nil
		},
	}

	httpReq := httptest.NewRequest("GET", getAuthorPath, nil)
	w := httptest.NewRecorder()

	handler := Get(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf(expectedGetAuthorStatusMsg, http.StatusCreated, w.Code)
	}
}

func TestGetMissingId(t *testing.T) {
	mockService := &mockGetService{
		getFunc: func(id string) (*model.Author, error) {
			return nil, nil
		},
	}

	httpReq := httptest.NewRequest("GET", "/author/get", nil)
	w := httptest.NewRecorder()

	handler := Get(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf(expectedGetAuthorStatusMsg, http.StatusBadRequest, w.Code)
	}
}

func TestGetServiceError(t *testing.T) {
	mockService := &mockGetService{
		getFunc: func(id string) (*model.Author, error) {
			return nil, errors.New("not found")
		},
	}

	httpReq := httptest.NewRequest("GET", getAuthorPath, nil)
	w := httptest.NewRecorder()

	handler := Get(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf(expectedGetAuthorStatusMsg, http.StatusNotFound, w.Code)
	}
}

func TestGetNilAuthor(t *testing.T) {
	mockService := &mockGetService{
		getFunc: func(id string) (*model.Author, error) {
			return nil, nil
		},
	}

	httpReq := httptest.NewRequest("GET", getAuthorPath, nil)
	w := httptest.NewRecorder()

	handler := Get(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf(expectedGetAuthorStatusMsg, http.StatusNotFound, w.Code)
}
}
