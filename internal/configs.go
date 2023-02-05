package metrics

type ServerConfig struct {
	Address       string
	Restore       bool
	StoreInterval int
	StoreFile     string
	Key           string
}

func checkServerConfig(envs *ServerEnvConfig, clies *ServerCLIOptions) *ServerConfig {
	var addr string = clies.Address
	var storeint string = clies.StoreInterval
	var storefile string = clies.StoreFile
	var cliRest string = clies.Restore
	var rest bool
	var storeintNumeric int = int(clies.GetNumericInterval("StoreInterval"))
	var key string = clies.Key
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
	if envs.Key != "" && envs.Key != key {
		key = envs.Key
	}
	return &ServerConfig{
		Address:       addr,
		StoreInterval: storeintNumeric,
		StoreFile:     storefile,
		Restore:       rest,
		Key:           key,
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
	Key            string
}

func checkAgentConfig(envs *AgentEnvConfig, clies *AgentCLIOptions) *AgentConfig {
	var addr string = clies.Address
	var pollint string = clies.PollInterval
	var reportint string = clies.ReportInterval
	var key string = clies.Key
	var reportintNumeric int = int(clies.GetNumericInterval("ReportInterval"))
	var pollintNumeric int = int(clies.GetNumericInterval("PollInterval"))
	if envs.Address != ADDRESS && envs.Address != addr {
		addr = envs.Address
	}
	if envs.ReportInterval != REPORT_INTERVAL && envs.ReportInterval != "" && envs.ReportInterval != reportint {
		reportintNumeric = int(envs.GetNumericInterval("ReportInterval"))
	}
	if envs.PollInterval != POLL_INTERVAL && envs.PollInterval != "" && envs.PollInterval != pollint {
		pollintNumeric = int(envs.GetNumericInterval("PollInterval"))
	}
	if envs.Key != "" && envs.Key != key {
		key = envs.Key
	}
	return &AgentConfig{
		Address:        addr,
		PollInterval:   pollintNumeric,
		ReportInterval: reportintNumeric,
		Key:            key,
	}
}

func GetAgentConfig() *AgentConfig {
	envConfig, err := NewAgentEnvConfig()
	if err != nil {
		InfoLog.Println("could not get env for server config, getting data from cli...")
		cliConfig := NewAgentCliOptions()
		return checkAgentConfig(envConfig, cliConfig)
	}
	return &AgentConfig{
		Address:        envConfig.Address,
		PollInterval:   int(envConfig.GetNumericInterval("PollInterval")),
		ReportInterval: int(envConfig.GetNumericInterval("ReportInterval")),
		Key:            envConfig.Key,
	}
}
