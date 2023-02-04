package metrics

type ServerConfig struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
}

func GetServerConfig(config *EnvConfig, args *ServerCLIOptions) *ServerConfig {
	serverConfig := ServerConfig{}
	if config.Address == "127.0.0.1:8080" {
		serverConfig.Address = args.Address
	} else {
		serverConfig.Address = config.Address
	}
	if config.Restore {
		serverConfig.Restore = config.Restore
	} else {
		serverConfig.Restore = args.Restore
	}
	if config.StoreFile != "/tmp/devops-metrics-db.json" {
		serverConfig.StoreFile = config.StoreFile
	} else {
		serverConfig.StoreFile = args.StoreFile
	}
	if config.StoreInterval == "300s" {
		serverConfig.StoreInterval = func() int {
			store, err := args.GetNumericInterval("StoreInterval")
			if err != nil {
				ErrorLog.Printf("error getting StoreInterval from CLI args: %e", err)
			}
			return int(store)
		}()
	} else {
		serverConfig.StoreInterval = func() int {
			store, err := config.GetNumericInterval("StoreInterval")
			if err != nil {
				ErrorLog.Printf("error getting StoreInterval from Env args: %e", err)
			}
			return int(store)
		}()
	}
	return &serverConfig
}

type AgentConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func GetAgentConfig(config *EnvConfig, args *AgentCLIOptions) *AgentConfig {
	agentConfig := AgentConfig{}
	if config.Address == "127.0.0.1:8080" {
		agentConfig.Address = args.Address
	} else {
		agentConfig.Address = config.Address
	}
	if config.PollInterval == "" {
		agentConfig.PollInterval = func() int {
			poll, err := args.GetNumericInterval("PollInterval")
			if err != nil {
				ErrorLog.Printf("error getting PollInterval from CLI args: %e", err)
			}
			return int(poll)
		}()
	} else {
		agentConfig.PollInterval = func() int {
			poll, err := config.GetNumericInterval("PollInterval")
			if err != nil {
				ErrorLog.Printf("error getting PollInterval from Env args: %e", err)
			}
			return int(poll)
		}()
	}
	if config.ReportInterval == "" {
		agentConfig.ReportInterval = func() int {
			rep, err := args.GetNumericInterval("ReportInterval")
			if err != nil {
				ErrorLog.Printf("error getting ReportInterval from CLI args: %e", err)
			}
			return int(rep)
		}()
	} else {
		agentConfig.ReportInterval = func() int {
			rep, err := config.GetNumericInterval("ReportInterval")
			if err != nil {
				ErrorLog.Printf("error getting ReportInterval from Env args: %e", err)
			}
			return int(rep)
		}()
	}
	return &agentConfig
}
