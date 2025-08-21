package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"coversion-de-monedas/internal/core"
	"coversion-de-monedas/internal/usecase"
	"github.com/sirupsen/logrus"
)

// ErrorResponse representa una respuesta de error estandarizada
type ErrorResponse struct {
	Error string `json:"error"`
}

// CurrencyHandler maneja las solicitudes HTTP relacionadas con el cambio de moneda
type CurrencyHandler struct {
	uc     *usecase.CurrencyConverterUseCase
	logger *logrus.Logger
}

// NewCurrencyHandler crea una nueva instancia de CurrencyHandler
func NewCurrencyHandler(uc *usecase.CurrencyConverterUseCase, logger *logrus.Logger) *CurrencyHandler {
	return &CurrencyHandler{
		uc:     uc,
		logger: logger,
	}
}

// GetExchangeRate maneja las solicitudes para obtener la tasa de cambio actual
func (h *CurrencyHandler) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	// Obtener el parámetro de moneda de la URL (opcional, por defecto USD)
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "USD"
	}

	h.logger.Infof("Obteniendo tasa de cambio para %s", currency)
	rate, err := h.uc.GetExchangeRate(currency)
	if err != nil {
		h.logger.Errorf("Error al obtener la tasa de cambio para %s: %v", currency, err)
		respondWithError(w, http.StatusInternalServerError, "Error al obtener la tasa de cambio")
		return
	}

	respondWithJSON(w, http.StatusOK, rate, h.logger)
}

// GetAllRates maneja las solicitudes para obtener todas las tasas de cambio
func (h *CurrencyHandler) GetAllRates(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Obteniendo todas las tasas de cambio")
	rates, err := h.uc.GetAllRates()
	if err != nil {
		h.logger.Errorf("Error al obtener las tasas de cambio: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al obtener las tasas de cambio")
		return
	}

	respondWithJSON(w, http.StatusOK, rates, h.logger)
}

// ConvertRequest representa la estructura de la solicitud de conversión
type ConvertRequest struct {
	Amount   float64 `json:"amount"`   // Cantidad a convertir
	From     string  `json:"from"`     // Moneda de origen ("USD", "EUR" o "VES")
	To       string  `json:"to"`       // Moneda de destino ("USD", "EUR" o "VES")
}

// ConvertCurrency maneja las solicitudes de conversión de moneda
func (h *CurrencyHandler) ConvertCurrency(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Procesando solicitud de conversión")
	
	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warnf("Solicitud inválida: %v", err)
		respondWithError(w, http.StatusBadRequest, "Solicitud inválida")
		return
	}

	// Validar y normalizar monedas
	req.From = normalizeCurrency(req.From)
	req.To = normalizeCurrency(req.To)

	if req.From == "" || req.To == "" {
		respondWithError(w, http.StatusBadRequest, "Debe especificar monedas de origen y destino válidas")
		return
	}

	if req.From == req.To {
		// Si las monedas son iguales, devolver la misma cantidad
		respondWithJSON(w, http.StatusOK, &core.ConversionResult{
			FromAmount: req.Amount,
			ToAmount:   req.Amount,
			Rate:       1,
			FromSymbol: req.From,
			ToSymbol:   req.To,
		}, h.logger)
		return
	}

	// Realizar la conversión
	result, err := h.uc.Convert(req.Amount, req.From, req.To)
	if err != nil {
		h.logger.Errorf("Error al realizar la conversión: %v", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Infof("Conversión exitosa: %f %s -> %f %s", req.Amount, req.From, result.ToAmount, req.To)
	respondWithJSON(w, http.StatusOK, result, h.logger)
}

// normalizeCurrency normaliza el formato de la moneda (mayúsculas, sin espacios)
func normalizeCurrency(currency string) string {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "USD", "$", "US$", "US":
		return "USD"
	case "EUR", "€", "EURO", "EUROS":
		return "EUR"
	case "VES", "BS", "BS.", "Bs", "Bs.", "bs", "bs.", "BOLIVAR", "BOLIVARES":
		return "VES"
	default:
		return "" // Moneda no reconocida
	}
}

// respondWithError envía una respuesta de error en formato JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message}, nil)
}

// respondWithJSON envía una respuesta en formato JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}, logger *logrus.Logger) {
	response, err := json.Marshal(payload)
	if err != nil {
		if logger != nil {
			logger.Errorf("Error al serializar la respuesta: %v", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error interno del servidor"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}
