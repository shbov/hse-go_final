package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shbov/hse-go_final/internal/location/docs"
	"github.com/shbov/hse-go_final/internal/location/model/requests"
	"github.com/shbov/hse-go_final/internal/location/service"
	"github.com/toshi0607/chi-prometheus"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"log"

	"github.com/juju/zaputil/zapctx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"moul.io/chizap"
)

var (
	ServiceName = "location-service"
	tracer      = otel.Tracer(ServiceName)
)

type adapter struct {
	config  *Config
	service service.Location
	server  *http.Server
}

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
// @Param        lat    body     requests.SetDriverLocationBody  true  "Latitude and longitude  in decimal degrees"
// @Success 200
// @Failure 400
// @Router /{driver_id}/location [post]
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

	writeResponse(w, http.StatusOK, "Success operation")
}

// GetDriversByLocation godoc
// @Summary GetDriversByLocation
// @Description Поиск водителей по заданным координатам и радиусу
// @Accept       json
// @Success 200
// @Router / [get]
func (a *adapter) GetDriversByLocation(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "SetDriverLocation")
	defer span.End()

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var request requests.GetDriversByLocationReqBody
	err := decoder.Decode(&request)
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

	writeJSONResponse(w, http.StatusOK, result)
}

func (a *adapter) Serve(ctx context.Context) error {
	lg := zapctx.Logger(ctx)

	shut := initTracerProvider(a.config.OtlpAddress)
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

	apiRouter.Post("/{driver_id}/location", a.SetDriverLocation)
	apiRouter.Get("/", a.GetDriversByLocation)

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
	_ = a.server.Shutdown(ctx)
}

func New(
	config *Config,
	service service.Location) Adapter {

	if config.SwaggerAddress != "" {
		docs.SwaggerInfo.Host = config.SwaggerAddress
	} else {
		docs.SwaggerInfo.Host = config.ServeAddress
	}

	docs.SwaggerInfo.BasePath = config.BasePath

	log.Println("adapter successfully created")
	return &adapter{
		config:  config,
		service: service,
	}
}

func initTracerProvider(otlpAddress string) func() {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(ServiceName),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("otlp address is " + otlpAddress)
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpAddress),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	sctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	traceExp, err := otlptrace.New(sctx, traceClient)
	if err != nil {
		log.Fatal(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
