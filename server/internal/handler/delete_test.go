package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteByID(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		deleteFunc func(id string) error
		wantStatus int
	}{
		{
			name: "success",
			path: "/resource?id=1",
			deleteFunc: func(id string) error {
				if id != "1" {
					t.Fatalf("unexpected id: %s", id)
				}
				return nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing id",
			path:       "/resource",
			deleteFunc: func(string) error { return nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "service error",
			path:       "/resource?id=1",
			deleteFunc: func(string) error { return errors.New("boom") },
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tc.path, nil)
			w := httptest.NewRecorder()

			DeleteByID(tc.deleteFunc)(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, w.Code)
			}
		})
	}
}
