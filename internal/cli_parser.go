package metrics

import (
	"flag"
	"strconv"
	"strings"
)

type ServerCLIOptions struct {
	Address       string
	Restore       string
	StoreInterval string
	StoreFile     string
	Key           string
	Database      string
}

type AgentCLIOptions struct {
	Address        string
	ReportInterval string
	PollInterval   string
	Key            string
}

func (scli *ServerCLIOptions) GetNumericInterval(intervalName string) int64 {
	switch intervalName {
	case "StoreInterval":
		multiplier := getMultiplier(scli.StoreInterval)
		stringValue := strings.Split(scli.StoreInterval, scli.StoreInterval[len(scli.StoreInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value
	}

	return 0
}

func (acli *AgentCLIOptions) GetNumericInterval(intervalName string) int64 {
	switch intervalName {
	case "ReportInterval":
		multiplier := getMultiplier(acli.ReportInterval)
		stringValue := strings.Split(acli.ReportInterval, acli.ReportInterval[len(acli.ReportInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value
	case "PollInterval":
		multiplier := getMultiplier(acli.PollInterval)
		stringValue := strings.Split(acli.PollInterval, acli.PollInterval[len(acli.PollInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return *multiplier * value
	}

	return 0
}

func NewServerCliOptions() *ServerCLIOptions {
	var address = flag.String("a", "localhost:8080", "server address")
	var restore = flag.String("r", "default", "restore metrics from file")
	var storeInterval = flag.String("i", "30s", "period between file save")
	var storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	var key = flag.String("k", "", "key to calculate hash")
	var database = flag.String("d", "", "database link")
	flag.Parse()

	return &ServerCLIOptions{
		Address:       *address,
		Restore:       *restore,
		StoreInterval: *storeInterval,
		StoreFile:     *storeFile,
		Key:           *key,
		Database:      *database,
	}
}

func NewAgentCliOptions() *AgentCLIOptions {
	var address = flag.String("a", "localhost:8080", "server address")
	var reportInterval = flag.String("r", "10s", "report interval time")
	var pollInterval = flag.String("p", "2s", "poll interval time")
	var key = flag.String("k", "", "key to calculate hash")
	flag.Parse()

	return &AgentCLIOptions{
		Address:        *address,
		ReportInterval: *reportInterval,
		PollInterval:   *pollInterval,
		Key:            *key,
	}
}
