package envparser

import (
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"

	"github.com/nmramorov/gowatcher/internal/errors"
)

const (
	Address        string = "127.0.0.1:8080"
	ReportInterval string = "10s"
	PollInterval   string = "2s"
	StoreInterval  string = "300s"
	StoreFile      string = "/tmp/devops-metrics-db.json"
	Restore        string = "default"
	Key            string = ""
	DatabaseDSN    string = ""
	RateLimit      int    = 0
)

type AgentEnvConfig struct {
	Address        string `env:"ADDRESS,required"`
	ReportInterval string `env:"REPORT_INTERVAL,required"`
	PollInterval   string `env:"POLL_INTERVAL,required"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY"`
	Config         string `env:"CONFIG"`
	GRPC           bool `env:"GRPC"`
}

func checkAgentEnvs(envs *AgentEnvConfig) *AgentEnvConfig {
	addr := Address
	pollint := PollInterval
	reportint := ReportInterval
	key := Key
	rate := RateLimit
	if envs.Address != Address && envs.Address != "" {
		addr = envs.Address
	}
	if envs.ReportInterval != ReportInterval && envs.ReportInterval != "" {
		reportint = envs.ReportInterval
	}
	if envs.PollInterval != PollInterval && envs.PollInterval != "" {
		pollint = envs.PollInterval
	}
	if envs.Key != "" {
		key = envs.Key
	}
	if envs.RateLimit != 0 {
		rate = envs.RateLimit
	}
	return &AgentEnvConfig{
		Address:        addr,
		PollInterval:   pollint,
		ReportInterval: reportint,
		Key:            key,
		RateLimit:      rate,
		CryptoKey:      envs.CryptoKey,
		Config:         envs.Config,
		GRPC:           envs.GRPC,
	}
}

func NewAgentEnvConfig() (*AgentEnvConfig, error) {
	var config AgentEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return checkAgentEnvs(&config), errors.ErrorWithEnvConfig
	}
	return &config, nil
}

type ServerEnvConfig struct {
	Address       string `env:"ADDRESS,required"`
	StoreInterval string `env:"STORE_INTERVAL,required"`
	StoreFile     string `env:"STORE_FILE,required"`
	Restore       string `env:"RESTORE,required"`
	Key           string `env:"KEY"`
	Database      string `env:"DATABASE_DSN"`
	CryptoKey     string `env:"CRYPTO_KEY"`
	Config        string `env:"CONFIG"`
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
	GRPC          bool `env:"GRPC"`
}

func checkServerEnvs(envs *ServerEnvConfig) *ServerEnvConfig {
	addr := Address
	storeint := StoreInterval
	storefile := StoreFile
	rest := Restore
	key := Key
	db := DatabaseDSN
	if envs.Address != Address && envs.Address != "" {
		addr = envs.Address
	}
	if envs.Restore != "default" && envs.Restore != "" {
		rest = envs.Restore
	}
	if envs.StoreFile != StoreFile && envs.StoreFile != "" {
		storefile = envs.StoreFile
	}
	if envs.StoreInterval != StoreInterval && envs.StoreInterval != "" {
		storeint = envs.StoreInterval
	}
	if envs.Key != "" {
		key = envs.Key
	}
	if envs.Database != "" {
		db = envs.Database
	}
	return &ServerEnvConfig{
		Address:       addr,
		StoreInterval: storeint,
		StoreFile:     storefile,
		Restore:       rest,
		Key:           key,
		Database:      db,
		CryptoKey:     envs.CryptoKey,
		Config:        envs.Config,
		TrustedSubnet: envs.TrustedSubnet,
		GRPC:          envs.GRPC,
	}
}

func NewServerEnvConfig() (*ServerEnvConfig, error) {
	var config ServerEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return checkServerEnvs(&config), errors.ErrorWithEnvConfig
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
	}
	return &multiplier
}

func (e *AgentEnvConfig) GetNumericInterval(intervalName string) int64 {
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
	}
	return 0
}

func (e *ServerEnvConfig) GetNumericInterval(intervalName string) int64 {
	if intervalName == "StoreInterval" {
		multiplier := getMultiplier(e.StoreInterval)
		stringValue := strings.Split(e.StoreInterval, e.StoreInterval[len(e.StoreInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value
	}

	return 0
}
