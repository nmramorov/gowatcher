package cliparser

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/nmramorov/gowatcher/internal/errors"
	"github.com/nmramorov/gowatcher/internal/log"
)

type ServerCLIOptions struct {
	Address       string
	Restore       string
	StoreInterval string
	StoreFile     string
	Key           string
	Database      string
	CryptoKey     string
	Config        string
	TrustedSubnet string
	GRPC          bool
}

type AgentCLIOptions struct {
	Address        string
	ReportInterval string
	PollInterval   string
	Key            string
	RateLimit      int
	CryptoKey      string
	Config         string
	GRPC           bool
}

func (scli *ServerCLIOptions) GetNumericInterval(intervalName string) int64 {
	if intervalName == "StoreInterval" {
		multiplier := GetMultiplier(scli.StoreInterval)
		stringValue := strings.Split(scli.StoreInterval, scli.StoreInterval[len(scli.StoreInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return multiplier * value
	}

	return 0
}

func (acli *AgentCLIOptions) GetNumericInterval(intervalName string) int64 {
	switch intervalName {
	case "ReportInterval":
		multiplier := GetMultiplier(acli.ReportInterval)
		stringValue := strings.Split(acli.ReportInterval, acli.ReportInterval[len(acli.ReportInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return multiplier * value
	case "PollInterval":
		multiplier := GetMultiplier(acli.PollInterval)
		stringValue := strings.Split(acli.PollInterval, acli.PollInterval[len(acli.PollInterval)-1:])[0]
		value, _ := strconv.ParseInt(stringValue, 10, 64)
		return multiplier * value
	}

	return 0
}

func NewServerCliOptions() (*ServerCLIOptions, error) {
	serverOptions := flag.NewFlagSet("server options", flag.ContinueOnError)
	address := serverOptions.String("a", "localhost:8080", "server address")
	restore := serverOptions.String("r", "default", "restore metrics from file")
	storeInterval := serverOptions.String("i", "30s", "period between file save")
	storeFile := serverOptions.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	key := serverOptions.String("k", "", "key to calculate hash")
	database := serverOptions.String("d", "", "database link")
	cryptoKey := serverOptions.String("crypto-key", "", "path to private key")
	config := serverOptions.String("c", "", "server json config path")
	subnet := serverOptions.String("t", "", "trusted subnet CIDR")
	grpc := serverOptions.Bool("grpc", false, "grpc mode")
	if err := serverOptions.Parse(os.Args[1:]); err != nil {
		log.ErrorLog.Printf("error parsing server cli options: %e", err)
		return nil, errors.ErrorWithCli
	}

	return &ServerCLIOptions{
		Address:       *address,
		Restore:       *restore,
		StoreInterval: *storeInterval,
		StoreFile:     *storeFile,
		Key:           *key,
		Database:      *database,
		CryptoKey:     *cryptoKey,
		Config:        *config,
		TrustedSubnet: *subnet,
		GRPC:          *grpc,
	}, nil
}

func NewAgentCliOptions() (*AgentCLIOptions, error) {
	agentOptions := flag.NewFlagSet("agent options", flag.ContinueOnError)
	address := agentOptions.String("a", "localhost:8080", "server address")
	reportInterval := agentOptions.String("r", "10s", "report interval time")
	pollInterval := agentOptions.String("p", "2s", "poll interval time")
	key := agentOptions.String("k", "", "key to calculate hash")
	rate := agentOptions.Int("l", 0, "rate limit")
	cryptoKey := agentOptions.String("crypto-key", "", "path to public key")
	config := agentOptions.String("c", "", "path to agent json config")
	grpc := agentOptions.Bool("grpc", false, "grpc mode")
	if err := agentOptions.Parse(os.Args[1:]); err != nil {
		log.ErrorLog.Printf("error parsing agent cli options: %e", err)
		return nil, errors.ErrorWithCli
	}

	return &AgentCLIOptions{
		Address:        *address,
		ReportInterval: *reportInterval,
		PollInterval:   *pollInterval,
		Key:            *key,
		RateLimit:      *rate,
		CryptoKey:      *cryptoKey,
		Config:         *config,
		GRPC:           *grpc,
	}, nil
}

func GetMultiplier(intervalValue string) int64 {
	var multiplier int64
	splitter := intervalValue[len(intervalValue)-1:]

	switch splitter {
	case `s`:
		multiplier = 1
	case `m`:
		multiplier = 60
	}
	return multiplier
}
