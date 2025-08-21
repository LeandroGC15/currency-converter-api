package config

import (
	"os"
	"time"
)

type ServerConfig struct {
	Port          string        `env:"PORT" envDefault:"8080"` // Sin los dos puntos
	AllowedOrigins []string     `env:"ALLOWED_ORIGINS" envDefault:"*"`
	EnableCORS    bool          `env:"ENABLE_CORS" envDefault:"true"`
	ReadTimeout   time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout  time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout   time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	CacheTTL      time.Duration `env:"CACHE_TTL" envDefault:"5m"`
}

func NewServerConfig() *ServerConfig {
	// Obtener el puerto de la variable de entorno o usar 8080 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"  // Sin los dos puntos aquí
	}

	// Configurar timeouts
	readTimeout, _ := time.ParseDuration("10s")
	writeTimeout, _ := time.ParseDuration("15s")
	idleTimeout, _ := time.ParseDuration("60s")
	cacheTTL, _ := time.ParseDuration("5m")

	// Obtener CACHE_TTL de las variables de entorno si está definido
	if ttl := os.Getenv("CACHE_TTL"); ttl != "" {
		if d, err := time.ParseDuration(ttl); err == nil {
			cacheTTL = d
		}
	}

	return &ServerConfig{
		Port:          port,  // Sin los dos puntos
		AllowedOrigins: []string{"*"},
		EnableCORS:    true,
		ReadTimeout:   readTimeout,
		WriteTimeout:  writeTimeout,
		IdleTimeout:   idleTimeout,
		CacheTTL:      cacheTTL,
	}
}
