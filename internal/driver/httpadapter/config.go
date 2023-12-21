package httpadapter

type Config struct {
	ServeAddress string `yaml:"serve_address"`
	BasePath     string `yaml:"base_path"`

	AccessTokenCookie  string `yaml:"access_token_cookie"`
	RefreshTokenCookie string `yaml:"refresh_token_cookie"`

	SwaggerAddress string `yaml:"swagger_address"`

	OtlpAddress string `yaml:"otlp"`
}
