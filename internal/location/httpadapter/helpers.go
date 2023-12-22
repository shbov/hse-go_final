package httpadapter

import (
	"github.com/go-errors/errors"
	"github.com/shbov/hse-go_final/internal/location/service"
	"github.com/shbov/hse-go_final/pkg/httpHelpers"
	"net/http"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrRequestIsIncorrect):
		httpHelpers.WriteJSONResponse(w, http.StatusBadRequest, httpHelpers.Error{Message: err.Error()})
	default:
		httpHelpers.WriteJSONResponse(w, http.StatusInternalServerError, httpHelpers.Error{Message: err.Error()})
	}
}
