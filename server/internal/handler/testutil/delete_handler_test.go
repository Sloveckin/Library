package testutil

import (
	"net/http"
	"testing"
)

func TestRunDeleteHandlerContract(t *testing.T) {
	RunDeleteHandlerContract(
		t,
		func(deleteFunc func(id string) error) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				id := r.URL.Query().Get("id")
				if id == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if err := deleteFunc(id); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusCreated)
			}
		},
		"/resource",
		"expected status %d, got %d",
	)
}
