package httpadapter

import (
	"github.com/go-errors/errors"
	"github.com/shbov/hse-go_final/internal/location/service"
	"github.com/shbov/hse-go_final/pkg/http_helpers"
	"net/http"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrRequestIsIncorrect):
		http_helpers.WriteJSONResponse(w, http.StatusBadRequest, http_helpers.Error{Message: err.Error()})
	default:
		http_helpers.WriteJSONResponse(w, http.StatusInternalServerError, http_helpers.Error{Message: err.Error()})
	}
}
