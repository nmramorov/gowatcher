package metrics

import (
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

type EnvConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval string `env:"REPORT_INTERVAL"`
	PollInterval   string `env:"POLL_INTERVAL"`
	StoreInterval  string `env:"STORE_INTERVAL"`
	StoreFile      string `env:"STORE_FILE"`
	Restore        bool   `env:"RESTORE"`
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

func (e *EnvConfig) GetNumericInterval(intervalName string) (int64, error) {
	switch intervalName {
	case "ReportInterval":
		multiplier := getMultiplier(e.ReportInterval)
		stringValue := strings.Split(e.ReportInterval, e.ReportInterval[len(e.ReportInterval)-1:])[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value, err
	case "PollInterval":
		multiplier := getMultiplier(e.PollInterval)
		stringValue := strings.Split(e.PollInterval, e.PollInterval[len(e.PollInterval)-1:])[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value, err
	case "StoreInterval":
		multiplier := getMultiplier(e.StoreInterval)
		stringValue := strings.Split(e.StoreInterval, e.StoreInterval[len(e.StoreInterval)-1:])[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value, err
	}

	return 0, ErrorWithIntervalConvertion
}
