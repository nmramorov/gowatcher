package metrics

import (
	"flag"
)

type ServerConfig struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
}

func checkServerConfig(envs *ServerEnvConfig, clies *ServerCLIOptions) *ServerConfig {
	var addr string = clies.Address
	var storeint string = clies.StoreInterval
	var storefile string = clies.StoreFile
	var cliRest string = clies.Restore
	var rest bool
	var storeintNumeric int = int(clies.GetNumericInterval("StoreInterval"))
	if envs.Address != ADDRESS && envs.Address != addr {
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
	if envs.StoreFile != STORE_FILE && envs.StoreFile != storefile {
		storefile = envs.StoreFile
	}
	if envs.StoreInterval != STORE_INTERVAL && envs.StoreInterval != "" && envs.StoreInterval != storeint {
		storeintNumeric = int(envs.GetNumericInterval("StoreInterval"))
	}
	return &ServerConfig{
		Address:       addr,
		StoreInterval: storeintNumeric,
		StoreFile:     storefile,
		Restore:       rest,
	}
}

func GetServerConfig() *ServerConfig {
	envConfig, err := NewServerEnvConfig()
	InfoLog.Println(envConfig, err)
	if err != nil {
		InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig := NewServerCliOptions()
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
	}
}

type AgentConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func GetAgentConfig() *AgentConfig {
	envConfig, err := NewAgentEnvConfig()
	if err != nil {
		InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig := NewAgentCliOptions()
		if flag.NFlag() == 3 {
			return &AgentConfig{
				Address:        cliConfig.Address,
				PollInterval:   int(cliConfig.GetNumericInterval("PollInterval")),
				ReportInterval: int(cliConfig.GetNumericInterval("ReportInterval")),
			}
		}
	}
	return &AgentConfig{
		Address:        envConfig.Address,
		PollInterval:   int(envConfig.GetNumericInterval("PollInterval")),
		ReportInterval: int(envConfig.GetNumericInterval("ReportInterval")),
	}
}
