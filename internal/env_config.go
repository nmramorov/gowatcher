package metrics

import (
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	Address        string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval string `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   string `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  string `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile      string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore        bool   `env:"RESTORE" envDefault:"true"`
}

func NewConfig() (*EnvConfig, error) {
	var config EnvConfig
	err := env.Parse(&config)
	if err != nil {
		return &config, ErrorWithEnvConfig
	}
	return &config, nil
}

func getMultiplier(intervalValue string) *int64 {
	var multiplier int64
	splitter := intervalValue[len(intervalValue)-1:]

	switch splitter {
	case `s`:
		multiplier = 1
	case `m`:
		multiplier = 60
	default:
		multiplier = 1
	}
	return &multiplier
}

func (e *EnvConfig) GetNumericInterval(intervalName string) int64 {
	switch intervalName {
	case "ReportInterval":
		multiplier := getMultiplier(e.ReportInterval)
		stringValue := strings.Split(e.ReportInterval, e.ReportInterval[len(e.ReportInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value
	case "PollInterval":
		multiplier := getMultiplier(e.PollInterval)
		stringValue := strings.Split(e.PollInterval, e.PollInterval[len(e.PollInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value
	case "StoreInterval":
		multiplier := getMultiplier(e.StoreInterval)
		stringValue := strings.Split(e.StoreInterval, e.StoreInterval[len(e.StoreInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value
	}

	return 0
}
