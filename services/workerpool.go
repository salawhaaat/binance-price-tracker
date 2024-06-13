package services

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	binance "github.com/binance/binance-connector-go"
	"github.com/salawhaaat/binance-price-tracker/config"
	"github.com/salawhaaat/binance-price-tracker/models"
)

type WorkerPool struct {
	MaxWorkers   int
	RequestCount int
	Symbols      []string
	ResultBuffer *bytes.Buffer // Buffer to store result strings
	Client       *binance.Client
	Prices       sync.Map
	mu           sync.Mutex // Mutex for ResultBuffer
	countMu      sync.Mutex // Mutex for RequestCount
	wg           sync.WaitGroup
}

func NewWorkerPool(cfg config.Config, client *binance.Client) *WorkerPool {
	return &WorkerPool{
		MaxWorkers:   cfg.MaxWorkers,
		Symbols:      cfg.Symbols,
		Client:       client,
		RequestCount: 0,
		ResultBuffer: bytes.NewBufferString(""),
		mu:           sync.Mutex{},
		countMu:      sync.Mutex{},
	}
}

// Run starts the worker pool and listens for stop signal
func (wp *WorkerPool) Run(stopChan chan struct{}) {
	wp.wg.Add(wp.MaxWorkers)
	for i := 0; i < wp.MaxWorkers; i++ {
		go func(workerId int) {
			defer wp.wg.Done()
			for {
				select {
				case <-stopChan:
					time.Sleep(5 * time.Second) // Wait for all workers to finish
					return
				default:
					j := wp.GetRequestsCount() % len(wp.Symbols)
					wp.Worker(&wp.Symbols[j])
				}
			}
		}(i)
	}
	wp.wg.Wait()
}

// Worker is a function that makes a GET request to the Binance API
func (wp *WorkerPool) Worker(symbol *string) {
	r := prepareRequest(symbol)

	resp, err := wp.Client.HTTPClient.Get(r)
	if err != nil {
		log.Println("Error making GET request for ", *symbol, err)
		return
	}
	defer resp.Body.Close()

	var ticker models.SymbolPrice
	if err := json.NewDecoder(resp.Body).Decode(&ticker); err != nil {
		log.Println("Error decoding response for ", *symbol, err)
		return
	}

	res := ticker.String()
	if wp.PriceChange(&ticker) {
		res += " changed"
	}
	wp.WriteToBuffer(res + "\n")
	wp.IncrementRequestCount()
}
