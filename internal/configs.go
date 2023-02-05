package metrics

import (
	"flag"
	"fmt"
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
		InfoLog.Println(flag.NFlag())
		// if flag.NFlag() == 4 {
		// 	return &ServerConfig{
		// 		Restore:       cliConfig.Restore,
		// 		Address:       cliConfig.Address,
		// 		StoreInterval: int(cliConfig.GetNumericInterval("StoreInterval")),
		// 		StoreFile:     cliConfig.StoreFile,
		// 	}
		// }
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
	// serverConfig := ServerConfig{}
	// if config.Address == "127.0.0.1:8080" {
	// 	serverConfig.Address = args.Address
	// } else {
	// 	serverConfig.Address = config.Address
	// }
	// if config.Restore {
	// 	serverConfig.Restore = config.Restore
	// } else {
	// 	serverConfig.Restore = args.Restore
	// }
	// if config.StoreFile != "/tmp/devops-metrics-db.json" {
	// 	serverConfig.StoreFile = config.StoreFile
	// } else {
	// 	serverConfig.StoreFile = args.StoreFile
	// }
	// if config.StoreInterval == "300s" {
	// 	serverConfig.StoreInterval = func() int {
	// 		store, err := args.GetNumericInterval("StoreInterval")
	// 		if err != nil {
	// 			ErrorLog.Printf("error getting StoreInterval from CLI args: %e", err)
	// 		}
	// 		return int(store)
	// 	}()
	// } else {
	// 	serverConfig.StoreInterval = func() int {
	// 		store, err := config.GetNumericInterval("StoreInterval")
	// 		if err != nil {
	// 			ErrorLog.Printf("error getting StoreInterval from Env args: %e", err)
	// 		}
	// 		return int(store)
	// 	}()
	// }
	// return &serverConfig
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
		fmt.Println(flag.NFlag())
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
	// agentConfig := AgentConfig{}
	// if config.Address == "127.0.0.1:8080" {
	// 	agentConfig.Address = args.Address
	// } else {
	// 	agentConfig.Address = config.Address
	// }
	// if config.PollInterval == "" {
	// 	agentConfig.PollInterval = func() int {
	// 		poll, err := args.GetNumericInterval("PollInterval")
	// 		if err != nil {
	// 			ErrorLog.Printf("error getting PollInterval from CLI args: %e", err)
	// 		}
	// 		return int(poll)
	// 	}()
	// } else {
	// 	agentConfig.PollInterval = func() int {
	// 		poll, err := config.GetNumericInterval("PollInterval")
	// 		if err != nil {
	// 			ErrorLog.Printf("error getting PollInterval from Env args: %e", err)
	// 		}
	// 		return int(poll)
	// 	}()
	// }
	// if config.ReportInterval == "" {
	// 	agentConfig.ReportInterval = func() int {
	// 		rep, err := args.GetNumericInterval("ReportInterval")
	// 		if err != nil {
	// 			ErrorLog.Printf("error getting ReportInterval from CLI args: %e", err)
	// 		}
	// 		return int(rep)
	// 	}()
	// } else {
	// 	agentConfig.ReportInterval = func() int {
	// 		rep, err := config.GetNumericInterval("ReportInterval")
	// 		if err != nil {
	// 			ErrorLog.Printf("error getting ReportInterval from Env args: %e", err)
	// 		}
	// 		return int(rep)
	// 	}()
	// }
	// return &agentConfig
}
