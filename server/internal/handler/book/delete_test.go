package book

import (
	"net/http"
	"server/internal/handler/testutil"
	"testing"
)

const expectedDeleteBookStatusMsg = "Expected status %d, got %d"

type mockDeleteService struct {
	deleteFunc func(id string) error
}

func (m *mockDeleteService) Delete(id string) error {
	return m.deleteFunc(id)
}

func TestDeleteContract(t *testing.T) {
	testutil.RunDeleteHandlerContract(
		t,
		func(deleteFunc func(id string) error) http.HandlerFunc {
			return Delete(&mockDeleteService{deleteFunc: deleteFunc})
		},
		"/book/delete",
		expectedDeleteBookStatusMsg,
	)
}
