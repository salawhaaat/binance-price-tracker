package config

import (
	"flag"
	"log"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Symbols    []string `yaml:"symbols"`
	MaxWorkers int      `yaml:"max_workers"`
}

var Cfg Config

var (
	ApiKey    string
	SecretKey string
)

// ParseFlags parses the command line flags
func ParseFlags() {
	flag.StringVar(&ApiKey, "api-key", "", "Binance API key")
	flag.StringVar(&SecretKey, "secret-key", "", "Binance API secret key")
	flag.Parse()
}

// InitConfig reads the configuration file and sets the configuration values
func InitConfig() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse() // -config=config/config.toml
	if err := cleanenv.ReadConfig(*configPath, &Cfg); err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}
	if MaxWorkers := runtime.NumCPU(); Cfg.MaxWorkers > MaxWorkers || Cfg.MaxWorkers <= 0 {
		Cfg.MaxWorkers = MaxWorkers // set to max value if invalid or not set
	}
}
