package metrics

import (
	"flag"
	"strconv"
	"strings"
)

type ServerCLIOptions struct {
	Address       string
	Restore       bool
	StoreInterval string
	StoreFile     string
}

type AgentCLIOptions struct {
	Address        string
	ReportInterval string
	PollInterval   string
}

func (scli *ServerCLIOptions) GetNumericInterval(intervalName string) (int64, error) {
	switch intervalName {
	case "StoreInterval":
		multiplier := getMultiplier(scli.StoreInterval)
		stringValue := strings.Split(scli.StoreInterval, scli.StoreInterval[len(scli.StoreInterval)-1:])[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value, err
	}

	return 0, ErrorWithIntervalConvertion
}

func (acli *AgentCLIOptions) GetNumericInterval(intervalName string) (int64, error) {
	switch intervalName {
	case "ReportInterval":
		multiplier := getMultiplier(acli.ReportInterval)
		stringValue := strings.Split(acli.ReportInterval, acli.ReportInterval[len(acli.ReportInterval)-1:])[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value, err
	case "PollInterval":
		multiplier := getMultiplier(acli.PollInterval)
		stringValue := strings.Split(acli.PollInterval, acli.PollInterval[len(acli.PollInterval)-1:])[0]
		value, err := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value, err
	}

	return 0, ErrorWithIntervalConvertion
}

func NewServerCliOptions() *ServerCLIOptions {
	var address = flag.String("a", "localhost:8080", "server address")
	var restore = flag.Bool("r", true, "restore metrics from file")
	var storeInterval = flag.String("i", "30s", "period between file save")
	var storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	flag.Parse()

	return &ServerCLIOptions{
		Address:       *address,
		Restore:       *restore,
		StoreInterval: *storeInterval,
		StoreFile:     *storeFile,
	}
}

func NewAgentCliOptions() *AgentCLIOptions {
	var address = flag.String("a", "localhost:8080", "server address")
	var reportInterval = flag.String("r", "10s", "report interval time")
	var pollInterval = flag.String("p", "2s", "poll interval time")
	flag.Parse()
	return &AgentCLIOptions{
		Address:        *address,
		ReportInterval: *reportInterval,
		PollInterval:   *pollInterval,
	}
}
