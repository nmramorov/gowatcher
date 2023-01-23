package metrics

import (
	"github.com/spf13/pflag"
)

type ServerCLIOptions struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
}

type AgentCLIOptions struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func getConfig() *EnvConfig {
	config, err := NewConfig()
	if err != nil {
		ErrorLog.Printf("Error during server arguments handling: %e", err)
		panic(err)
	}
	return config
}

func NewServerOptions() *ServerCLIOptions {
	config := getConfig()
	storeInterval, err := config.GetNumericInterval("StoreInterval")
	if err != nil {
		panic(err)
	}
	return &ServerCLIOptions{
		Address:       *pflag.String("a", config.Address, "server address"),
		Restore:       *pflag.Bool("r", config.Restore, "restore metrics from file"),
		StoreInterval: *pflag.Int("i", int(storeInterval), "period between file save"),
		StoreFile:     *pflag.String("f", config.StoreFile, "name of file where metrics stored"),
	}
}

func NewAgentOptions() *AgentCLIOptions {
	config := getConfig()
	pollInterval, err := config.GetNumericInterval("PollInterval")
	if err != nil {
		panic(err)
	}
	reportInterval, err := config.GetNumericInterval("ReportInterval")
	if err != nil {
		panic(err)
	}

	return &AgentCLIOptions{
		Address:        *pflag.String("a", config.Address, "server address"),
		PollInterval:   *pflag.Int("p", int(pollInterval), "metrics poll interval"),
		ReportInterval: *pflag.Int("r", int(reportInterval), "metrics report interval"),
	}
}
