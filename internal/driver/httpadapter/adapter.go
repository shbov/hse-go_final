package httpadapter

import (
	"context"
	"fmt"
	"github.com/shbov/hse-go_final/internal/driver/docs"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"github.com/shbov/hse-go_final/internal/driver/service"
	"github.com/shbov/hse-go_final/pkg/http_helpers"
	tracer2 "github.com/shbov/hse-go_final/pkg/tracer"
	"github.com/toshi0607/chi-prometheus"

	"github.com/juju/zaputil/zapctx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"moul.io/chizap"
)

var _ Adapter = (*adapter)(nil)

var (
	ServiceName = "driver-service"
	tracer      = otel.Tracer(ServiceName)
)

type adapter struct {
	config       *Config
	kafkaService service.KafkaService
	tripService  service.Trip
	server       *http.Server
}

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

	// TODO: we need to link driver to a trip
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

	http_helpers.WriteResponse(w, http.StatusOK, "Successful operation")
}

func (a *adapter) Serve(ctx context.Context) error {
	lg := zapctx.Logger(ctx)

	shut := tracer2.InitTracerProvider(ctx, a.config.OtlpAddress, ServiceName)
	defer shut()

	r := chi.NewRouter()
	apiRouter := chi.NewRouter()

	apiRouter.Use(otelchi.Middleware(ServiceName, otelchi.WithChiRoutes(r)))

	m := chiprometheus.New("main")
	m.MustRegisterDefault()
	apiRouter.Use(m.Handler)

	apiRouter.Use(chizap.New(lg, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	apiRouter.Get("/trips", a.GetTrips)
	apiRouter.Get("/trips/{trip_id}", a.GetTripByTripId)

	apiRouter.Post("/trips/{trip_id}/cancel", a.CancelTrip)
	apiRouter.Post("/trips/{trip_id}/accept", a.AcceptTrip)
	apiRouter.Post("/trips/{trip_id}/start", a.StartTrip)
	apiRouter.Post("/trips/{trip_id}/end", a.EndTrip)

	// установка маршрута для документации
	// Адрес, по которому будет доступен doc.json
	apiRouter.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", a.config.BasePath))))

	r.Mount(a.config.BasePath, apiRouter)

	a.server = &http.Server{Addr: a.config.ServeAddress, Handler: r}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9000", nil)
		if err != nil {
			lg.Error(err.Error())
		}
	}()

	return a.server.ListenAndServe()
}

func (a *adapter) Shutdown(ctx context.Context) {
	lg := zapctx.Logger(ctx)
	err := a.server.Shutdown(ctx)
	if err != nil {
		lg.Error(err.Error())
	}
}

func New(
	ctx context.Context,
	config *Config,
	messageQueue service.KafkaService,
	tripRepo service.Trip) Adapter {
	lg := zapctx.Logger(ctx)

	if config.SwaggerAddress != "" {
		docs.SwaggerInfo.Host = config.SwaggerAddress
	} else {
		docs.SwaggerInfo.Host = config.ServeAddress
	}

	docs.SwaggerInfo.BasePath = config.BasePath

	lg.Info("adapter successfully created")
	return &adapter{
		config:       config,
		kafkaService: messageQueue,
		tripService:  tripRepo,
	}
}

type requestData struct {
	UserId string
	TripId string
	Trip   *trip.Trip
}

func (a *adapter) validate(r *http.Request) (*requestData, error) {
	userId := r.Header.Get("user_id")
	tripId := chi.URLParam(r, "trip_id")
	if !http_helpers.IsValidUUID(tripId) {
		return nil, fmt.Errorf("invalid uuid")
	}

	trip, err := a.tripService.GetTripByUserIdTripId(r.Context(), userId, tripId)
	if err != nil {
		return nil, fmt.Errorf("trip not found")
	}

	return &requestData{UserId: userId, TripId: tripId, Trip: trip}, nil
}
