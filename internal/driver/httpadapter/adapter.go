package httpadapter

import (
	"context"
	"fmt"
	"github.com/shbov/hse-go_final/internal/driver/docs"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
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
