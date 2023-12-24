package httpadapter

import (
	"context"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shbov/hse-go_final/internal/location/docs"
	"github.com/shbov/hse-go_final/internal/location/service"
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

	apiRouter.Get("/drivers", a.GetDriversByLocation)
	apiRouter.Post("/drivers/{driver_id}/location", a.SetDriverLocation)

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
