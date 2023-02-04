package metrics

import "flag"

type ServerConfig struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
}

func GetServerConfig() *ServerConfig {
	envConfig, err := NewConfig()
	InfoLog.Println(envConfig, err)
	if err != nil {
		InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig := NewServerCliOptions()
		InfoLog.Println(flag.Parsed())
		if flag.Parsed() {
			return &ServerConfig{
				Restore:       cliConfig.Restore,
				Address:       cliConfig.Address,
				StoreInterval: int(cliConfig.GetNumericInterval("StoreInterval")),
				StoreFile:     cliConfig.StoreFile,
			}
		}
	}
	return &ServerConfig{
		Restore:       envConfig.Restore,
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
	envConfig, err := NewConfig()
	if err != nil {
		InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig := NewServerCliOptions()
		return &AgentConfig{
			Address:        cliConfig.Address,
			PollInterval:   int(cliConfig.GetNumericInterval("PollInterval")),
			ReportInterval: int(cliConfig.GetNumericInterval("ReportInterval")),
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