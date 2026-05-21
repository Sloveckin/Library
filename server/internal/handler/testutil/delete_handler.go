package testutil

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DeleteHandlerFactory func(deleteFunc func(id string) error) http.HandlerFunc

func RunDeleteHandlerContract(
	t *testing.T,
	makeHandler DeleteHandlerFactory,
	deletePath string,
	expectedStatusMsg string,
) {
	t.Helper()

	t.Run("success", func(t *testing.T) {
		handler := makeHandler(func(string) error { return nil })
		req := httptest.NewRequest(http.MethodDelete, deletePath+"?id=1", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf(expectedStatusMsg, http.StatusCreated, w.Code)
		}
	})

	t.Run("missing id", func(t *testing.T) {
		handler := makeHandler(func(string) error { return nil })
		req := httptest.NewRequest(http.MethodDelete, deletePath, nil)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf(expectedStatusMsg, http.StatusBadRequest, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		handler := makeHandler(func(string) error { return errors.New("not found") })
		req := httptest.NewRequest(http.MethodDelete, deletePath+"?id=1", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf(expectedStatusMsg, http.StatusBadRequest, w.Code)
		}
	})
}
