package config

import (
	cli "github.com/nmramorov/gowatcher/internal/config/cli_parser"
	env "github.com/nmramorov/gowatcher/internal/config/env_parser"
	"github.com/nmramorov/gowatcher/internal/log"
)

var DEFAULT = "default"

type ServerConfig struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
	Key           string
	Database      string
}

func checkServerConfig(envs *env.ServerEnvConfig, clies *cli.ServerCLIOptions) *ServerConfig {
	addr := clies.Address
	storeint := clies.StoreInterval
	storefile := clies.StoreFile
	cliRest := clies.Restore
	rest := false
	storeintNumeric := int(clies.GetNumericInterval("StoreInterval"))
	key := clies.Key
	db := clies.Database
	if envs.Address != env.Address && envs.Address != addr {
		addr = envs.Address
	}
	if envs.Restore != env.Restore && envs.Restore != "" {
		if envs.Restore == "true" {
			rest = true
		} else {
			rest = false
		}
	}
	if envs.Restore == env.Restore && cliRest != DEFAULT {
		if cliRest == "true" {
			rest = true
		} else {
			rest = false
		}
	}
	if envs.StoreFile != env.StoreFile && envs.StoreFile != storefile {
		storefile = envs.StoreFile
	}
	if envs.StoreInterval != env.StoreInterval && envs.StoreInterval != "" && envs.StoreInterval != storeint {
		storeintNumeric = int(envs.GetNumericInterval("StoreInterval"))
	}
	if envs.Key != "" && envs.Key != key {
		key = envs.Key
	}
	if envs.Database != "" && envs.Database != db {
		db = envs.Database
	}
	return &ServerConfig{
		Address:       addr,
		StoreInterval: storeintNumeric,
		StoreFile:     storefile,
		Restore:       rest,
		Key:           key,
		Database:      db,
	}
}

func GetServerConfig() (*ServerConfig, error) {
	envConfig, err := env.NewServerEnvConfig()
	log.InfoLog.Println(envConfig, err)
	if err != nil {
		log.InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig, err := cli.NewServerCliOptions()
		if err != nil {
			return nil, err
		}
		return checkServerConfig(envConfig, cliConfig), nil
	}
	var rest bool
	if envConfig.Restore == "true" {
		rest = true
	} else {
		rest = false
	}
	return &ServerConfig{
		Restore:       rest,
		Address:       envConfig.Address,
		StoreInterval: int(envConfig.GetNumericInterval("StoreInterval")),
		StoreFile:     envConfig.StoreFile,
		Database:      envConfig.Database,
	}, nil
}

type AgentConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
	Key            string
	RateLimit      int
}

func checkAgentConfig(envs *env.AgentEnvConfig, clies *cli.AgentCLIOptions) *AgentConfig {
	addr := clies.Address
	pollint := clies.PollInterval
	reportint := clies.ReportInterval
	key := clies.Key
	reportintNumeric := int(clies.GetNumericInterval("ReportInterval"))
	pollintNumeric := int(clies.GetNumericInterval("PollInterval"))
	rate := clies.RateLimit
	if envs.Address != env.Address && envs.Address != addr {
		addr = envs.Address
	}
	if envs.ReportInterval != env.ReportInterval && envs.ReportInterval != "" && envs.ReportInterval != reportint {
		reportintNumeric = int(envs.GetNumericInterval("ReportInterval"))
	}
	if envs.PollInterval != env.PollInterval && envs.PollInterval != "" && envs.PollInterval != pollint {
		pollintNumeric = int(envs.GetNumericInterval("PollInterval"))
	}
	if envs.Key != "" && envs.Key != key {
		key = envs.Key
	}
	if envs.RateLimit != 0 && envs.RateLimit != rate {
		rate = envs.RateLimit
	}
	return &AgentConfig{
		Address:        addr,
		PollInterval:   pollintNumeric,
		ReportInterval: reportintNumeric,
		Key:            key,
		RateLimit:      rate,
	}
}

func GetAgentConfig() (*AgentConfig, error) {
	envConfig, err := env.NewAgentEnvConfig()
	if err != nil {
		log.InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig, err := cli.NewAgentCliOptions()
		if err != nil {
			return nil, err
		}
		return checkAgentConfig(envConfig, cliConfig), nil
	}
	return &AgentConfig{
		Address:        envConfig.Address,
		PollInterval:   int(envConfig.GetNumericInterval("PollInterval")),
		ReportInterval: int(envConfig.GetNumericInterval("ReportInterval")),
		Key:            envConfig.Key,
	}, nil
}
