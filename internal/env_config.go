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

func (e *EnvConfig) GetNumericInterval(intervalName string) (int64, error) {
	switch intervalName {
	case "ReportInterval":
		stringValue := strings.Split(e.ReportInterval, `s`)[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return value, err
	case "PollInterval":
		stringValue := strings.Split(e.PollInterval, `s`)[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return value, err
	case "StoreInterval":
		stringValue := strings.Split(e.StoreInterval, `s`)[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return value, err
	}

	return 0, ErrorWithIntervalConvertion
}
