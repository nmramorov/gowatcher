package metrics

import (
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

const (
	ADDRESS         string = "127.0.0.1:8080"
	REPORT_INTERVAL string = "10s"
	POLL_INTERVAL   string = "2s"
	STORE_INTERVAL  string = "300s"
	STORE_FILE      string = "/tmp/devops-metrics-db.json"
	RESTORE         string = "default"
	KEY             string = ""
	DATABASE_DSN    string = ""
	RATE_LIMIT      int    = 0
)

type AgentEnvConfig struct {
	Address        string `env:"ADDRESS,required"`
	ReportInterval string `env:"REPORT_INTERVAL,required"`
	PollInterval   string `env:"POLL_INTERVAL,required"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func checkAgentEnvs(envs *AgentEnvConfig) *AgentEnvConfig {
	var addr string = ADDRESS
	var pollint string = POLL_INTERVAL
	var reportint string = REPORT_INTERVAL
	var key string = KEY
	var rate int = RATE_LIMIT
	if envs.Address != ADDRESS && envs.Address != "" {
		addr = envs.Address
	}
	if envs.ReportInterval != REPORT_INTERVAL && envs.ReportInterval != "" {
		reportint = envs.ReportInterval
	}
	if envs.PollInterval != POLL_INTERVAL && envs.PollInterval != "" {
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
	}
}

func NewAgentEnvConfig() (*AgentEnvConfig, error) {
	var config AgentEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return checkAgentEnvs(&config), ErrorWithEnvConfig
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
}

func checkServerEnvs(envs *ServerEnvConfig) *ServerEnvConfig {
	var addr string = ADDRESS
	var storeint string = STORE_INTERVAL
	var storefile string = STORE_FILE
	var rest string = RESTORE
	var key string = KEY
	var db string = DATABASE_DSN
	if envs.Address != ADDRESS && envs.Address != "" {
		addr = envs.Address
	}
	if envs.Restore != "default" && envs.Restore != "" {
		rest = envs.Restore
	}
	if envs.StoreFile != STORE_FILE && envs.StoreFile != "" {
		storefile = envs.StoreFile
	}
	if envs.StoreInterval != STORE_INTERVAL && envs.StoreInterval != "" {
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
	}
}

func NewServerEnvConfig() (*ServerEnvConfig, error) {
	var config ServerEnvConfig
	err := env.Parse(&config)
	if err != nil {
		return checkServerEnvs(&config), ErrorWithEnvConfig
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
	switch intervalName {
	case "StoreInterval":
		multiplier := getMultiplier(e.StoreInterval)
		stringValue := strings.Split(e.StoreInterval, e.StoreInterval[len(e.StoreInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value
	}

	return 0
}
