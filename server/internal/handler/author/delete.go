package author

import (
	"net/http"
	v "server/internal/handler"
)

type deleteService interface {
	Delete(id string) error
}

func Delete(service deleteService) http.HandlerFunc {
	return v.DeleteByID(service.Delete)
}
