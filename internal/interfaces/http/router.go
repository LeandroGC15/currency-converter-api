package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Router configura y devuelve un enrutador HTTP
func NewRouter(handler *CurrencyHandler, logger *logrus.Logger) http.Handler {
	r := mux.NewRouter()

	// Configurar rutas de la API
	api := r.PathPrefix("/api/v1").Subrouter()
	
	// Health check
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Rutas de cambio de moneda
	api.HandleFunc("/rates", handler.GetExchangeRate).Methods("GET") // Obtener tasa para una moneda específica
	api.HandleFunc("/rates/all", handler.GetAllRates).Methods("GET") // Obtener todas las tasas disponibles
	api.HandleFunc("/convert", handler.ConvertCurrency).Methods("POST") // Realizar conversión

	// Configurar CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Añadir middleware de logging
	loggedRouter := addLoggingMiddleware(r, logger)

	// Aplicar CORS al router con logging
	return corsHandler.Handler(loggedRouter)
}

// addLoggingMiddleware añade logging a las peticiones HTTP
func addLoggingMiddleware(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
