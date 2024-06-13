package services

import "github.com/salawhaaat/binance-price-tracker/models"

func (wp *WorkerPool) WriteToBuffer(s string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	wp.ResultBuffer.WriteString(s)
}

func (wp *WorkerPool) PriceChange(t *models.SymbolPrice) bool {
	value, loaded := wp.Prices.LoadOrStore(t.Symbol, t.Price)
	if !loaded {
		return false
	}
	return value != t.Price
}

func (wp *WorkerPool) IncrementRequestCount() {
	wp.countMu.Lock()
	wp.RequestCount++
	wp.countMu.Unlock()
}

func (wp *WorkerPool) GetRequestsCount() int {
	wp.countMu.Lock()
	defer wp.countMu.Unlock()
	return wp.RequestCount
}

func (wp *WorkerPool) WaitForCompletion() {
	wp.wg.Wait()
}
