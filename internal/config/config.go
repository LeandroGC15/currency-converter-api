package config

type ServerConfig struct {
	Port         string `env:"PORT" envDefault:":8080"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envDefault:"*"`
	EnableCORS   bool   `env:"ENABLE_CORS" envDefault:"true"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:         ":8080",
		AllowedOrigins: []string{"*"},
		EnableCORS:   true,
	}
}
