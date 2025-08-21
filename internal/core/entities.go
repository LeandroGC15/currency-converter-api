package core

// ExchangeRate representa la tasa de cambio actual
type ExchangeRate struct {
	Currency   string  `json:"currency"`   // USD o EUR
	Price      float64 `json:"price"`      // Precio en VES
	Change     float64 `json:"change"`     // Cambio en valor absoluto
	Percent    float64 `json:"percent"`    // Cambio porcentual
	Symbol     string  `json:"symbol"`     // ▲ o ▼
	LastUpdate string  `json:"last_update"` // Fecha de actualización
}

type ConversionRequest struct {
	Amount float64 `json:"amount"`
	FromCurrency string `json:"from_currency"`
	ToCurrency string `json:"to_currency"`
}

// ConversionResult representa el resultado de una conversión
type ConversionResult struct {
	FromAmount float64 `json:"from_amount"` // Cantidad original
	ToAmount   float64 `json:"to_amount"`   // Cantidad convertida
	Rate       float64 `json:"rate"`        // Tasa de cambio utilizada
	FromSymbol string  `json:"from_symbol"` // Símbolo de la moneda de origen
	ToSymbol   string  `json:"to_symbol"`   // Símbolo de la moneda de destino
}