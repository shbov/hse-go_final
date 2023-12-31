package httpadapter

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/shbov/hse-go_final/internal/location/model/requests"
	"github.com/shbov/hse-go_final/internal/location/service"
	"github.com/shbov/hse-go_final/pkg/http_helpers"
	"net/http"
)

// @title Location service
// @version 1.0
// @description This is a location service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// SetDriverLocation godoc
// @Summary SetDriverLocation
// @Description Обновление данных о позиции водителя
// @Accept       json
// @Param        driver_id    path     string  true  "ID of driver"  Format(uuid)
// @Param        request	body     requests.SetDriverLocationBody  true  "Latitude and longitude  in decimal degrees"
// @Success 200
// @Failure 400
// @Router /drivers/{driver_id}/location [post]
func (a *adapter) SetDriverLocation(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "SetDriverLocation")
	defer span.End()

	driverId := chi.URLParam(r, "driver_id")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var request requests.SetDriverLocationBody
	err := decoder.Decode(&request)
	if err != nil {
		writeError(w, service.ErrRequestIsIncorrect)
		return
	}

	if !request.Validate() {
		writeError(w, service.ErrRequestIsIncorrect)
		return
	}

	err = a.service.SetLocationByDriverId(r.Context(), driverId, request.Lat, request.Lng)
	if err != nil {
		writeError(w, err)
		return
	}

	http_helpers.WriteResponse(w, http.StatusOK, "Success operation")
}

// GetDriversByLocation godoc
// @Summary GetDriversByLocation
// @Description Поиск водителей по заданным координатам и радиусу
// @Accept       json
// @Param        lat    query     float64  true  "Latitude in decimal degrees"
// @Param        lng    query     float64  true  "Longitude in decimal degrees"
// @Param        radius    query     float64  true  "Radius in meters"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /drivers [get]
func (a *adapter) GetDriversByLocation(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GetDriversByLocation")
	defer span.End()

	lat, err := http_helpers.ParseQueryToFloat("lat", r.URL.Query())
	if err != nil {
		writeError(w, service.ErrRequestIsIncorrect)
		return
	}

	lng, err := http_helpers.ParseQueryToFloat("lng", r.URL.Query())
	if err != nil {
		writeError(w, service.ErrRequestIsIncorrect)
		return
	}

	radius, err := http_helpers.ParseQueryToFloat("radius", r.URL.Query())
	if err != nil || radius <= 0 {
		writeError(w, service.ErrRequestIsIncorrect)
		return
	}

	request := requests.GetDriversByLocationReqBody{
		Lat:    lat,
		Lng:    lng,
		Radius: radius,
	}

	if err != nil {
		writeError(w, service.ErrRequestIsIncorrect)
		return
	}

	if !request.Validate() {
		writeError(w, service.ErrRequestIsIncorrect)
		return
	}

	result, err := a.service.GetDriversInLocation(r.Context(), request.Lat, request.Lng, request.Radius)
	if err != nil {
		writeError(w, err)
		return
	}

	http_helpers.WriteJSONResponse(w, http.StatusOK, result)
}
