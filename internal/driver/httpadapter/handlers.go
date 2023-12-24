package httpadapter

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"github.com/shbov/hse-go_final/pkg/http_helpers"
	"net/http"
)

// @title Driver service
// @version 1.0
// @description This is a driver service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// GetTrips godoc
// @Summary GetTrips
// @Description Получение списка поездок пользователя
// @Accept       json
// @Param        user_id    header     string  true  "ID of user"  Format(uuid)
// @Success 200
// @Failure 400
// @Router /trips [get]
func (a *adapter) GetTrips(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GetTrips")
	defer span.End()

	userId := r.Header.Get("user_id")

	trips, err := a.tripService.GetTripsByUserId(r.Context(), userId)
	if err != nil {
		writeError(w, err)
		return
	}

	http_helpers.WriteJSONResponse(w, http.StatusOK, trips)
}

// GetTripByTripId godoc
// @Summary GetTripByTripId
// @Description Получение поездки по ID
// @Accept       json
// @Param        user_id    header     string  true  "ID of user"  Format(uuid)
// @Param        trip_id    path     string  true  "ID of trip"  Format(uuid)
// @Success 200
// @Failure 400
// @Failure 404
// @Router /trips/{trip_id} [get]
func (a *adapter) GetTripByTripId(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GetTripByTripId")
	defer span.End()

	data, err := a.validate(r)
	if err != nil {
		writeError(w, err)
		return
	}

	if err != nil {
		writeError(w, err)
		return
	}

	http_helpers.WriteJSONResponse(w, http.StatusOK, *data.Trip)
}

// AcceptTrip godoc
// @Summary AcceptTrip
// @Description Принятие поездки водителем
// @Accept       json
// @Param        user_id    header     string  true  "ID of user"  Format(uuid)
// @Param        trip_id    path     string  true  "ID of trip"  Format(uuid)
// @Success 200
// @Failure 400
// @Failure 404
// @Router /trips/{trip_id}/accept [post]
func (a *adapter) AcceptTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "AcceptTrip")
	defer span.End()

	userID := r.Header.Get("user_id")
	tripID := chi.URLParam(r, "trip_id")

	if !http_helpers.IsValidUUID(tripID) {
		writeError(w, fmt.Errorf("invalid uuid"))
		return
	}

	err := a.tripService.UpdateDriverIdByTripId(r.Context(), tripID, userID)
	if err != nil {
		writeError(w, err)
		return
	}

	// mongodb <- status, driver_id
	err = a.tripService.ChangeTripStatus(r.Context(), tripID, trip_status.ACCEPTED)
	if err != nil {
		writeError(w, err)
		return
	}

	// -> kafka
	err = a.kafkaService.AcceptTrip(r.Context(), userID, tripID)
	if err != nil {
		writeError(w, err)
		return
	}

	http_helpers.WriteResponse(w, http.StatusOK, "Successful operation")
}

// CancelTrip godoc
// @Summary CancelTrip
// @Description Отмена поездки водителем
// @Accept       json
// @Param        user_id    header     string  true  "ID of user"  Format(uuid)
// @Param        trip_id    path     string  true  "ID of trip"  Format(uuid)
// @Param        reason    query     string  false  "Reason of cancel"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /trips/{trip_id}/cancel [post]
func (a *adapter) CancelTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "CancelTrip")
	defer span.End()

	reason := r.URL.Query().Get("reason")
	data, err := a.validate(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.kafkaService.CancelTrip(r.Context(), data.TripId, reason)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.tripService.ChangeTripStatus(r.Context(), data.TripId, trip_status.CANCELED)
	if err != nil {
		writeError(w, err)
		return
	}

	http_helpers.WriteResponse(w, http.StatusOK, "Successful operation")
}

// StartTrip godoc
// @Summary StartTrip
// @Description Начало поездки водителем
// @Accept       json
// @Param        user_id    header     string  true  "ID of user"  Format(uuid)
// @Param        trip_id    path     string  true  "ID of trip"  Format(uuid)
// @Success 200
// @Failure 400
// @Failure 404
// @Router /trips/{trip_id}/start [post]
func (a *adapter) StartTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "StartTrip")
	defer span.End()

	data, err := a.validate(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.kafkaService.StartTrip(r.Context(), data.TripId)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.tripService.ChangeTripStatus(r.Context(), data.TripId, trip_status.STARTED)
	if err != nil {
		writeError(w, err)
		return
	}

	http_helpers.WriteResponse(w, http.StatusOK, "Success operation")
}

// EndTrip godoc
// @Summary EndTrip
// @Description Окончание поездки водителем
// @Accept       json
// @Param        user_id    header     string  true  "ID of user"  Format(uuid)
// @Param        trip_id    path     string  true  "ID of trip"  Format(uuid)
// @Success 200
// @Failure 400
// @Failure 404
// @Router /trips/{trip_id}/end [post]
func (a *adapter) EndTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "EndTrip")
	defer span.End()

	data, err := a.validate(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.kafkaService.EndTrip(r.Context(), data.TripId)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.tripService.ChangeTripStatus(r.Context(), data.TripId, trip_status.ENDED)
	if err != nil {
		writeError(w, err)
		return
	}

	http_helpers.WriteResponse(w, http.StatusOK, "Successful operation")
}
