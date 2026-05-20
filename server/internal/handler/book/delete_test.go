package book

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const expectedDeleteBookStatusMsg = "Expected status %d, got %d"

type mockDeleteService struct {
	deleteFunc func(id string) error
}

func (m *mockDeleteService) Delete(id string) error {
	return m.deleteFunc(id)
}

func TestDeleteSuccess(t *testing.T) {
	mockService := &mockDeleteService{
		deleteFunc: func(id string) error {
			return nil
		},
	}

	httpReq := httptest.NewRequest("DELETE", "/book/delete?id=1", nil)
	w := httptest.NewRecorder()

	handler := Delete(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf(expectedDeleteBookStatusMsg, http.StatusCreated, w.Code)
	}
}

func TestDeleteMissingId(t *testing.T) {
	mockService := &mockDeleteService{
		deleteFunc: func(id string) error {
			return nil
		},
	}

	httpReq := httptest.NewRequest("DELETE", "/book/delete", nil)
	w := httptest.NewRecorder()

	handler := Delete(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf(expectedDeleteBookStatusMsg, http.StatusBadRequest, w.Code)
	}
}

func TestDeleteServiceError(t *testing.T) {
	mockService := &mockDeleteService{
		deleteFunc: func(id string) error {
			return errors.New("not found")
		},
	}

	httpReq := httptest.NewRequest("DELETE", "/book/delete?id=1", nil)
	w := httptest.NewRecorder()

	handler := Delete(mockService)
	handler(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf(expectedDeleteBookStatusMsg, http.StatusBadRequest, w.Code)
}
}
