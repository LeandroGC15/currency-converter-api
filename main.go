package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"coversion-de-monedas/internal/config"
	"coversion-de-monedas/internal/infrastructure"
	httpHandler "coversion-de-monedas/internal/interfaces/http"
	"coversion-de-monedas/internal/usecase"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configuración del logger
	logger := configureLogger()
	logger.Info("Iniciando aplicación de conversión de monedas")

	// Configuración
	cfg := config.NewServerConfig()

	// Inicializar dependencias
	repo := infrastructure.NewPyDolarRepository()
	uc := usecase.NewCurrencyConverterUseCase(repo)
	handler := httpHandler.NewCurrencyHandler(uc, logger)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":" + cfg.Port,  // Agregar los dos puntos aquí
		Handler:      httpHandler.NewRouter(handler, logger),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Canal para manejar señales de terminación
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		logger.Infof("Servidor escuchando en http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar señal de terminación
	<-done
	logger.Info("Recibida señal de apagado, cerrando servidor...")

	// Configurar tiempo de espera para el cierre
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Error al apagar el servidor: %v", err)
	}

	logger.Info("Servidor detenido correctamente")
}

// configureLogger configura el logger de la aplicación
func configureLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Configurar nivel de log basado en una variable de entorno o usar INFO por defecto
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warnf("Nivel de log inválido '%s', usando 'info' por defecto", logLevel)
		level = logrus.InfoLevel
	}

	logger.SetLevel(level)
	return logger
}
