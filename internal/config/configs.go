package config

import (
	cli "github.com/nmramorov/gowatcher/internal/config/cli_parser"
	env "github.com/nmramorov/gowatcher/internal/config/env_parser"
	"github.com/nmramorov/gowatcher/internal/log"
)

type ServerConfig struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
	Key           string
	Database      string
}

func checkServerConfig(envs *env.ServerEnvConfig, clies *cli.ServerCLIOptions) *ServerConfig {
	var addr string = clies.Address
	var storeint string = clies.StoreInterval
	var storefile string = clies.StoreFile
	var cliRest string = clies.Restore
	var rest bool
	var storeintNumeric int = int(clies.GetNumericInterval("StoreInterval"))
	var key string = clies.Key
	var db string = clies.Database
	if envs.Address != env.ADDRESS && envs.Address != addr {
		addr = envs.Address
	}
	if envs.Restore != "default" && envs.Restore != "" {
		if envs.Restore == "true" {
			rest = true
		} else {
			rest = false
		}
	}
	if envs.Restore == "default" && cliRest != "default" {
		if cliRest == "true" {
			rest = true
		} else {
			rest = false
		}
	}
	if envs.StoreFile != env.STORE_FILE && envs.StoreFile != storefile {
		storefile = envs.StoreFile
	}
	if envs.StoreInterval != env.STORE_INTERVAL && envs.StoreInterval != "" && envs.StoreInterval != storeint {
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

func GetServerConfig() *ServerConfig {
	envConfig, err := env.NewServerEnvConfig()
	log.InfoLog.Println(envConfig, err)
	if err != nil {
		log.InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig := cli.NewServerCliOptions()
		return checkServerConfig(envConfig, cliConfig)
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
	}
}

type AgentConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
	Key            string
	RateLimit      int
}

func checkAgentConfig(envs *env.AgentEnvConfig, clies *cli.AgentCLIOptions) *AgentConfig {
	var addr string = clies.Address
	var pollint string = clies.PollInterval
	var reportint string = clies.ReportInterval
	var key string = clies.Key
	var reportintNumeric int = int(clies.GetNumericInterval("ReportInterval"))
	var pollintNumeric int = int(clies.GetNumericInterval("PollInterval"))
	var rate int = clies.RateLimit
	if envs.Address != env.ADDRESS && envs.Address != addr {
		addr = envs.Address
	}
	if envs.ReportInterval != env.REPORT_INTERVAL && envs.ReportInterval != "" && envs.ReportInterval != reportint {
		reportintNumeric = int(envs.GetNumericInterval("ReportInterval"))
	}
	if envs.PollInterval != env.POLL_INTERVAL && envs.PollInterval != "" && envs.PollInterval != pollint {
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

func GetAgentConfig() *AgentConfig {
	envConfig, err := env.NewAgentEnvConfig()
	if err != nil {
		log.InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig := cli.NewAgentCliOptions()
		return checkAgentConfig(envConfig, cliConfig)
	}
	return &AgentConfig{
		Address:        envConfig.Address,
		PollInterval:   int(envConfig.GetNumericInterval("PollInterval")),
		ReportInterval: int(envConfig.GetNumericInterval("ReportInterval")),
		Key:            envConfig.Key,
	}
}
