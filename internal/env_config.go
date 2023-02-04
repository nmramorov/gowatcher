package metrics

import (
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

const (
	ADDRESS         string = "127.0.0.1"
	REPORT_INTERVAL        = "10s"
	POLL_INTERVAL          = "2s"
	STORE_INTERVAL         = "300s"
	STORE_FILE             = "/tmp/devops-metrics-db.json"
	RESTORE         bool   = true
)

type EnvConfig struct {
	Address        string `env:"ADDRESS,required"`
	ReportInterval string `env:"REPORT_INTERVAL,required"`
	PollInterval   string `env:"POLL_INTERVAL,required"`
	StoreInterval  string `env:"STORE_INTERVAL,required"`
	StoreFile      string `env:"STORE_FILE,required"`
	Restore        bool   `env:"RESTORE,required"`
}

func NewConfig() (*EnvConfig, error) {
	var config EnvConfig
	err := env.Parse(&config)
	if err != nil {
		return &EnvConfig{
			Address:        ADDRESS,
			StoreInterval:  STORE_INTERVAL,
			StoreFile:      STORE_FILE,
			Restore:        RESTORE,
			ReportInterval: REPORT_INTERVAL,
			PollInterval:   POLL_INTERVAL,
		}, ErrorWithEnvConfig
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
