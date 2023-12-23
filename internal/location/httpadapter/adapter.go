package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shbov/hse-go_final/internal/location/docs"
	"github.com/shbov/hse-go_final/internal/location/model/requests"
	"github.com/shbov/hse-go_final/internal/location/service"
	"github.com/shbov/hse-go_final/pkg/http_helpers"
	tracer2 "github.com/shbov/hse-go_final/pkg/tracer"
	"github.com/toshi0607/chi-prometheus"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
	"moul.io/chizap"
	"net/http"
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
// @Failure 404
// @Router / [get]
func (a *adapter) GetDriversByLocation(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GetDriversByLocation")
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

	http_helpers.WriteJSONResponse(w, http.StatusOK, result)
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

	apiRouter.Get("/", a.GetDriversByLocation)
	apiRouter.Post("/{driver_id}/location", a.SetDriverLocation)

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
	ctx context.Context,
	config *Config,
	service service.Location) Adapter {
	lg := zapctx.Logger(ctx)
	if config.SwaggerAddress != "" {
		docs.SwaggerInfo.Host = config.SwaggerAddress
	} else {
		docs.SwaggerInfo.Host = config.ServeAddress
	}

	docs.SwaggerInfo.BasePath = config.BasePath

	lg.Info("adapter successfully created")
	return &adapter{
		config:  config,
		service: service,
	}
}
