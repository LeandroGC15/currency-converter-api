package core

// ExchangeRateRepository define el contrato para obtener las tasas de cambio
type ExchangeRateRepository interface {
	// GetExchangeRate obtiene la tasa de cambio para una moneda específica (USD o EUR)
	GetExchangeRate(currency string) (*ExchangeRate, error)
	
	// GetAllRates obtiene todas las tasas de cambio disponibles
	GetAllRates() (map[string]*ExchangeRate, error)
}

type CurrencyConverter interface {
	Convert(amount float64, fromCurrency string, toCurrency string) (*ConversionResult, error)
	GetExchangeRate(currency string) (*ExchangeRate, error)
}