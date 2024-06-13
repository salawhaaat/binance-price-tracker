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
	MaxWorkers   int             // Number of workers to run
	RequestCount int             // Number of requests made
	Symbols      []string        // List of symbols to track
	ResultBuffer *bytes.Buffer   // Buffer to store result strings
	Client       *binance.Client // Binance client
	Prices       sync.Map        // Map to store symbol prices and check for changes
	mu           sync.Mutex      // Mutex for ResultBuffer
	countMu      sync.Mutex      // Mutex for RequestCount
	wg           sync.WaitGroup
}

func NewWorkerPool(cfg config.Config, client *binance.Client) *WorkerPool {
	return &WorkerPool{
		MaxWorkers:   cfg.MaxWorkers,            // Number of workers to run
		Symbols:      cfg.Symbols,               // List of symbols to track
		Client:       client,                    // Binance client
		RequestCount: 0,                         // Number of requests made
		ResultBuffer: bytes.NewBufferString(""), // Buffer to store results
		mu:           sync.Mutex{},              // Mutex for ResultBuffer
		countMu:      sync.Mutex{},              // Mutex for RequestCount
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
					j := wp.GetRequestsCount() % len(wp.Symbols) // Get the next symbol to process
					wp.Worker(&wp.Symbols[j])
				}
			}
		}(i)
	}
	wp.wg.Wait()
}

// Worker is a function that makes a GET request to the Binance API
func (wp *WorkerPool) Worker(symbol *string) {
	r := prepareRequest(symbol) // Prepare the request URL
	wp.IncrementRequestCount()  // Increment the request count
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

	res := ticker.String()       // Format the response
	if wp.PriceChange(&ticker) { // Check if the price has changed
		res += " changed"
	}
	wp.WriteToBuffer(res + "\n") // Write the response to the buffer

}
