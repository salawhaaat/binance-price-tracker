package main

import (
	"github.com/salawhaaat/binance-price-tracker/config"
	"github.com/salawhaaat/binance-price-tracker/services"

	binance_connector "github.com/binance/binance-connector-go"

	"bufio"
	"os"
)

func init() {
	config.ParseFlags()
	config.InitConfig()
}

func main() {
	client := binance_connector.NewClient(config.ApiKey, config.SecretKey)
	wp := services.NewWorkerPool(config.Cfg, client)
	stopChan := make(chan struct{})

	go wp.Run(stopChan)
	go services.StartResultPrinter(wp)

	// Listen for STOP input from stdin
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
