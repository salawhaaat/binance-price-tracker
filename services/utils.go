package services

import "github.com/salawhaaat/binance-price-tracker/models"

// WriteToBuffer writes a string to the buffer
func (wp *WorkerPool) WriteToBuffer(s string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.ResultBuffer.WriteString(s)
}

// PriceChange checks if the price of a symbol has changed
func (wp *WorkerPool) PriceChange(t *models.SymbolPrice) bool {
	value, loaded := wp.Prices.LoadOrStore(t.Symbol, t.Price)
	if !loaded {
		return false
	}
	return value != t.Price
}

// incrementRequestCount increments the request count
func (wp *WorkerPool) IncrementRequestCount() {
	wp.countMu.Lock()
	wp.RequestCount++
	wp.countMu.Unlock()
}

// GetRequestsCount returns the request count of the worker pool
func (wp *WorkerPool) GetRequestsCount() int {
	wp.countMu.Lock()
	defer wp.countMu.Unlock()
	return wp.RequestCount
}

// WaitForCompletion waits for all workers to finish
func (wp *WorkerPool) WaitForCompletion() {
	wp.wg.Wait()
}
