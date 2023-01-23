package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerCLIDefaults(t *testing.T) {
	config, err := NewConfig()
	if err != nil {
		ErrorLog.Printf("Error during server arguments handling: %e", err)
		panic(err)
	}
	interval, err := config.GetNumericInterval("StoreInterval")
	if err != nil {
		panic(err)
	}
	serverOptions := NewServerOptions()
	assert.Equal(t, config.Address, serverOptions.Address)
	assert.Equal(t, config.Restore, serverOptions.Restore)
	assert.Equal(t, config.StoreFile, serverOptions.StoreFile)
	assert.Equal(t, int(interval), serverOptions.StoreInterval)
}

func TestAgentCLIDefaults(t *testing.T) {
	config, err := NewConfig()
	if err != nil {
		ErrorLog.Printf("Error during agent arguments handling: %e", err)
		panic(err)
	}
	pollInterval, err := config.GetNumericInterval("PollInterval")
	if err != nil {
		panic(err)
	}
	reportInterval, err := config.GetNumericInterval("ReportInterval")
	if err != nil {
		panic(err)
	}
	agentOptions := NewAgentOptions()
	assert.Equal(t, config.Address, agentOptions.Address)
	assert.Equal(t, int(reportInterval), agentOptions.ReportInterval)
	assert.Equal(t, int(pollInterval), agentOptions.PollInterval)
}
