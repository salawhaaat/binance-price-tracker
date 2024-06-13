package models

import "fmt"

// SymbolPrice represents the price of a symbol
type SymbolPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (s *SymbolPrice) String() string {
	return fmt.Sprintf("%s: %s", s.Symbol, s.Price)
}
