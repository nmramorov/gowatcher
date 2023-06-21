package cliparser

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
	RateLimit      int
}

func (scli *ServerCLIOptions) GetNumericInterval(intervalName string) int64 {
	if intervalName == "StoreInterval" {
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
	address := flag.String("a", "localhost:8080", "server address")
	restore := flag.String("r", "default", "restore metrics from file")
	storeInterval := flag.String("i", "30s", "period between file save")
	storeFile := flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	key := flag.String("k", "", "key to calculate hash")
	database := flag.String("d", "", "database link")
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
	address := flag.String("a", "localhost:8080", "server address")
	reportInterval := flag.String("r", "10s", "report interval time")
	pollInterval := flag.String("p", "2s", "poll interval time")
	key := flag.String("k", "", "key to calculate hash")
	rate := flag.Int("l", 0, "rate limit")
	flag.Parse()

	return &AgentCLIOptions{
		Address:        *address,
		ReportInterval: *reportInterval,
		PollInterval:   *pollInterval,
		Key:            *key,
		RateLimit:      *rate,
	}
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
