# Binance Price Tracker

A Go application for tracking cryptocurrency prices on Binance using a worker pool.

## Features

- Concurrent price tracking
- Configurable workers and symbols
- Real-time price change detection
- Simple start/stop mechanism

## Prerequisite

- Go 1.16+
- Binance API Key and Secret Key

## Installation

1. Clone the repo:

   ```sh
   git clone https://github.com/salawhaaat/binance-price-tracker.git
   cd binance-price-tracker
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Configuration

Create `config.yaml`:

```yaml
symbols:
  - BTCUSDT
  - ETHUSDT
max_workers: 4
```

## Usage

Run with your Binance API key and secret key:

```sh
go run main.go -api-key=your_api_key -secret-key=your_secret_key -config=config.yaml
```

To stop the application, type `STOP` and press Enter.

## Project Structure

- `main.go`: Entry point.
- `config/`: Configuration handling.
- `models/`: Data models.
- `services/`: Core functionality.

## Code Overview

### Main Function

Initializes configuration and starts worker pool and result printer.

```go
func main() {
    client := binance_connector.NewClient(config.ApiKey, config.SecretKey)
    wp := services.NewWorkerPool(config.Cfg, client)
    stopChan := make(chan struct{})

    go wp.Run(stopChan)
    go services.StartResultPrinter(wp)

    reader := bufio.NewReader(os.Stdin)
    for {
        input, _ := reader.ReadString('\n')
        if input == "STOP\n" {
            close(stopChan)
            break
        }
    }

    wp.WaitForCompletion()
}
```

### Configuration

The number of max_workers must not exceed the number of virtual processor cores and must be forcibly set to this value if it exceeds it

```go
func InitConfig() {
    ... // Reads command-line flags and config file.
	if MaxWorkers := runtime.NumCPU(); Cfg.MaxWorkers > MaxWorkers || Cfg.MaxWorkers <= 0 {
		Cfg.MaxWorkers = MaxWorkers
	}
}
```

### Worker Pool

Handles concurrent requests to the Binance API and listens for STOP command from stopChan.
Because there is 5 sec delivery between prints, all request that were send before will be successfully printed, but others after print wont be canceled, just not be printed.

```go
func (wp *WorkerPool) Run(stopChan chan struct{}) {
    wp.wg.Add(wp.MaxWorkers)
    for i := 0; i < wp.MaxWorkers; i++ {
        go func(workerId int) {
            defer wp.wg.Done()
            for {
                select {
                case <-stopChan:
                    time.Sleep(5 * time.Second)
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
```

### Limitations

workerpool.go

```go
    // First, perfomance issues will occur on big numbers
    j := wp.GetRequestsCount() % len(wp.Symbols)
    // Second, There is no gracefull shutdown, just delay for last request.
```
