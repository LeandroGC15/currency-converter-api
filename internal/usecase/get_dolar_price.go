package usecase

import (
	"fmt"
	"strings"

	"coversion-de-monedas/internal/core"
)

type CurrencyConverterUseCase struct {
	repo core.ExchangeRateRepository
}

func NewCurrencyConverterUseCase(repo core.ExchangeRateRepository) *CurrencyConverterUseCase {
	return &CurrencyConverterUseCase{
		repo: repo,
	}
}

// GetExchangeRate obtiene la tasa de cambio para una moneda específica
func (uc *CurrencyConverterUseCase) GetExchangeRate(currency string) (*core.ExchangeRate, error) {
	return uc.repo.GetExchangeRate(currency)
}

// GetAllRates obtiene todas las tasas de cambio disponibles
func (uc *CurrencyConverterUseCase) GetAllRates() (map[string]*core.ExchangeRate, error) {
	return uc.repo.GetAllRates()
}

// Convert realiza la conversión de moneda según los parámetros proporcionados
func (uc *CurrencyConverterUseCase) Convert(amount float64, fromCurrency, toCurrency string) (*core.ConversionResult, error) {
	// Validar monedas
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)

	if fromCurrency == toCurrency {
		return &core.ConversionResult{
			FromAmount: amount,
			ToAmount:   amount,
			Rate:       1,
			FromSymbol: fromCurrency,
			ToSymbol:   toCurrency,
		}, nil
	}

	// Obtener tasas necesarias
	rates, err := uc.repo.GetAllRates()
	if err != nil {
		return nil, fmt.Errorf("error al obtener las tasas de cambio: %v", err)
	}

	// Verificar si las monedas son soportadas
	fromRate, fromIsVES := rates[fromCurrency]
	toRate, toIsVES := rates[toCurrency]

	if !fromIsVES && fromCurrency != "VES" {
		return nil, fmt.Errorf("moneda de origen no soportada: %s", fromCurrency)
	}

	if !toIsVES && toCurrency != "VES" {
		return nil, fmt.Errorf("moneda de destino no soportada: %s", toCurrency)
	}

	result := &core.ConversionResult{
		FromAmount: amount,
		FromSymbol: fromCurrency,
		ToSymbol:   toCurrency,
	}

	// Realizar la conversión
	if fromCurrency == "VES" {
		// De VES a otra moneda
		result.ToAmount = amount / toRate.Price
		result.Rate = 1 / toRate.Price
	} else if toCurrency == "VES" {
		// De otra moneda a VES
		result.ToAmount = amount * fromRate.Price
		result.Rate = fromRate.Price
	} else {
		// Entre dos monedas que no son VES (ej: USD a EUR)
		// Primero a VES, luego a la moneda destino
		vesAmount := amount * fromRate.Price
		result.ToAmount = vesAmount / toRate.Price
		result.Rate = fromRate.Price / toRate.Price
	}

	// Redondear a 2 decimales para mejor presentación
	result.ToAmount = roundToTwoDecimals(result.ToAmount)

	return result, nil
}

// roundToTwoDecimals redondea un número a 2 decimales
func roundToTwoDecimals(num float64) float64 {
	return float64(int(num*100)) / 100
}
