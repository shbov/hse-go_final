package httpadapter

import (
	"context"
	"fmt"
	"github.com/shbov/hse-go_final/internal/driver/docs"
	"github.com/shbov/hse-go_final/internal/driver/model"
	"github.com/shbov/hse-go_final/internal/driver/service"
	"github.com/shbov/hse-go_final/pkg/httpHelpers"
	tracer2 "github.com/shbov/hse-go_final/pkg/tracer"
	"github.com/toshi0607/chi-prometheus"

	"log"

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
	messageQueue service.KafkaService
	tripRepo     service.Trip
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

func (a *adapter) GetTrips(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GetTrips")
	defer span.End()

	userId := r.Header.Get("user_id")

	trips, err := a.tripRepo.GetTripsByUserId(r.Context(), userId)
	if err != nil {
		writeError(w, err)
		return
	}

	httpHelpers.WriteJSONResponse(w, http.StatusOK, trips)
}

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

	httpHelpers.WriteJSONResponse(w, http.StatusOK, *data.Trip)
}

func (a *adapter) AcceptTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "AcceptTrip")
	defer span.End()

	data, err := a.validate(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.messageQueue.AcceptTrip(r.Context(), data.UserId, data.TripId)
	if err != nil {
		writeError(w, err)
		return
	}

	httpHelpers.WriteResponse(w, http.StatusOK, "Successful operation")
}

func (a *adapter) CancelTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "CancelTrip")
	defer span.End()

	reason := r.URL.Query().Get("reason")
	data, err := a.validate(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.messageQueue.CancelTrip(r.Context(), data.TripId, reason)
	if err != nil {
		writeError(w, err)
		return
	}

	httpHelpers.WriteResponse(w, http.StatusOK, "Successful operation")
}

func (a *adapter) StartTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "StartTrip")
	defer span.End()

	data, err := a.validate(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.messageQueue.StartTrip(r.Context(), data.TripId)
	if err != nil {
		writeError(w, err)
		return
	}

	httpHelpers.WriteResponse(w, http.StatusOK, "Success operation")
}

func (a *adapter) EndTrip(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "EndTrip")
	defer span.End()

	data, err := a.validate(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.messageQueue.EndTrip(r.Context(), data.TripId)
	if err != nil {
		writeError(w, err)
		return
	}

	httpHelpers.WriteResponse(w, http.StatusOK, "Successful operation")
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

	// установка маршрута для документации
	// Адрес, по которому будет доступен doc.json
	apiRouter.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", a.config.BasePath))))

	r.Mount(a.config.BasePath, apiRouter)

	a.server = &http.Server{Addr: a.config.ServeAddress, Handler: r}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9000", nil)
		if err != nil {
			lg.Fatal(err.Error())
		}
	}()

	return a.server.ListenAndServe()
}

func (a *adapter) Shutdown(ctx context.Context) {
	lg := zapctx.Logger(ctx)
	err := a.server.Shutdown(ctx)
	if err != nil {
		lg.Fatal(err.Error())
	}
}

func New(
	config *Config,
	messageQueue service.KafkaService,
	tripRepo service.Trip) Adapter {

	if config.SwaggerAddress != "" {
		docs.SwaggerInfo.Host = config.SwaggerAddress
	} else {
		docs.SwaggerInfo.Host = config.ServeAddress
	}

	docs.SwaggerInfo.BasePath = config.BasePath

	log.Println("adapter successfully created")
	return &adapter{
		config:       config,
		messageQueue: messageQueue,
		tripRepo:     tripRepo,
	}
}

type requestData struct {
	UserId string
	TripId string
	Trip   *model.Trip
}

func (a *adapter) validate(r *http.Request) (*requestData, error) {
	userId := r.Header.Get("user_id")
	tripId := chi.URLParam(r, "trip_id")
	if !httpHelpers.IsValidUUID(tripId) {
		return nil, fmt.Errorf("invalid uuid")
	}

	trip, err := a.tripRepo.GetTripByUserIdTripId(r.Context(), userId, tripId)
	if err != nil {
		return nil, fmt.Errorf("trip not found")
	}

	return &requestData{UserId: userId, TripId: tripId, Trip: trip}, nil
}
