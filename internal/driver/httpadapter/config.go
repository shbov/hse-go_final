package httpadapter

type Config struct {
	ServeAddress string `yaml:"serve_address"`
	BasePath     string `yaml:"base_path"`

	SwaggerAddress string `yaml:"swagger_address"`

	OtlpAddress string `yaml:"otlp"`
}
