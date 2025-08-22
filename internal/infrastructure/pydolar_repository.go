package infrastructure

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"coversion-de-monedas/internal/core"
	"coversion-de-monedas/internal/infrastructure/cache"
)

const (
	cacheTTL = 5 * time.Minute
)

type PyDolarRepository struct {
	cache     *cache.TTLCache
	cacheLock sync.RWMutex
}

func NewPyDolarRepository() *PyDolarRepository {
	return &PyDolarRepository{
		cache: cache.NewTTLCache(),
	}
}

// GetExchangeRate obtiene la tasa de cambio para una moneda específica (USD o EUR)
func (r *PyDolarRepository) GetExchangeRate(currency string) (*core.ExchangeRate, error) {
	// Normalizar el código de moneda
	currency = strings.ToUpper(currency)
	
	// Validar moneda
	if currency != "USD" && currency != "EUR" {
		return nil, fmt.Errorf("moneda no soportada: %s. Use 'USD' o 'EUR'", currency)
	}

	r.cacheLock.RLock()
	cached, found := r.cache.Get(currency)
	r.cacheLock.RUnlock()

	if found {
		if rate, ok := cached.(*core.ExchangeRate); ok {
			return rate, nil
		}
	}

	// Si no está en caché o es inválido, obtener la tasa específica
	cmd := exec.Command("python3", "get_dolar.py", "--currency", strings.ToLower(currency))
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("error al ejecutar el script Python: %v, stderr: %s", err, string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("error al ejecutar el script Python: %v", err)
	}

	// Parsear la respuesta JSON del script Python
	var response struct {
		Success bool                `json:"success"`
		Data    *core.ExchangeRate `json:"data"`
		Error   string             `json:"error,omitempty"`
	}

	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta para %s: %v", currency, err)
	}

	if !response.Success {
		return nil, fmt.Errorf("error al obtener la tasa para %s: %s", currency, response.Error)
	}

	if response.Data == nil {
		return nil, fmt.Errorf("no se recibieron datos para la moneda %s", currency)
	}

	// Asegurarse de que la moneda esté establecida correctamente
	response.Data.Currency = currency

	// Actualizar caché
	r.cacheLock.Lock()
	r.cache.Set(currency, response.Data, cacheTTL)
	r.cacheLock.Unlock()

	return response.Data, nil
}

// GetAllRates obtiene todas las tasas de cambio disponibles
func (r *PyDolarRepository) GetAllRates() (map[string]*core.ExchangeRate, error) {
	r.cacheLock.RLock()
	cached, found := r.cache.Get("all_rates")
	r.cacheLock.RUnlock()

	if found {
		if rates, ok := cached.(map[string]*core.ExchangeRate); ok {
			return rates, nil
		}
	}

	// Si no está en caché o es inválido, obtener de la API
	return r.fetchAllRates()
}

// fetchAllRates obtiene todas las tasas de la API
func (r *PyDolarRepository) fetchAllRates() (map[string]*core.ExchangeRate, error) {
	rates := make(map[string]*core.ExchangeRate)

	// Obtener tasas para USD y EUR
	for _, currency := range []string{"USD", "EUR"} {
		cmd := exec.Command("python3", "get_dolar.py", "--currency", strings.ToLower(currency))
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("error al obtener la tasa para %s: %v", currency, err)
		}

		// Parsear la respuesta JSON del script Python
		var response struct {
			Success bool                  `json:"success"`
			Data    *core.ExchangeRate   `json:"data"`
			Error   string               `json:"error,omitempty"`
		}

		if err := json.Unmarshal(output, &response); err != nil {
			return nil, fmt.Errorf("error al decodificar la respuesta para %s: %v", currency, err)
		}

		if !response.Success {
			return nil, fmt.Errorf("error al obtener la tasa para %s: %s", currency, response.Error)
		}

		if response.Data == nil {
			return nil, fmt.Errorf("no se recibieron datos para la moneda %s", currency)
		}

		// Asegurarse de que la moneda esté establecida correctamente
		response.Data.Currency = currency
		rates[currency] = response.Data

		// Actualizar caché individual con TTL
		r.cacheLock.Lock()
		r.cache.Set(currency, response.Data, cacheTTL)
		r.cacheLock.Unlock()
	}

	// Actualizar caché de todas las tasas con TTL
	r.cacheLock.Lock()
	r.cache.Set("all_rates", rates, cacheTTL)
	r.cacheLock.Unlock()

	return rates, nil
}
