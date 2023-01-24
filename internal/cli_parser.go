package metrics

import (
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
