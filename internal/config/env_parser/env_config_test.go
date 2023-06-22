package envparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvConfig(t *testing.T) {
	testConfig, _ := NewServerEnvConfig()
	assert.Equal(t, testConfig.Address, `127.0.0.1:8080`)
	assert.Equal(t, testConfig.StoreInterval, `300s`)
	assert.Equal(t, testConfig.StoreFile, "/tmp/devops-metrics-db.json")
	assert.Equal(t, testConfig.Restore, "default")
	assert.Equal(t, testConfig.Key, "")
	assert.Equal(t, testConfig.Database, "")
}

func TestEnvConfigIntervalConvertion(t *testing.T) {
	testConfig, _ := NewAgentEnvConfig()
	testConfig.PollInterval = "2s"
	testConfig.ReportInterval = "10m"
	reportInterval := testConfig.GetNumericInterval("ReportInterval")
	assert.Equal(t, reportInterval, int64(600))
	pollInterval := testConfig.GetNumericInterval("PollInterval")
	assert.Equal(t, pollInterval, int64(2))
	assert.Equal(t, testConfig.RateLimit, 0)
}

func TestGetNumericInterval(t *testing.T) {
	testConfig, _ := NewServerEnvConfig()
	testConfig.StoreInterval = "10s"
	interval := testConfig.GetNumericInterval("StoreInterval")
	assert.Equal(t, int64(10), interval)

	wrongIntervalNameRes := testConfig.GetNumericInterval("Some field")
	assert.Equal(t, int64(0), wrongIntervalNameRes)
}

func TestCheckAgentEnvs(t *testing.T) {
	config, _ := NewAgentEnvConfig()
	config.Address = "some addr"
	config.Key = "some key"
	config.PollInterval = "my interval"
	config.RateLimit = 222
	config.ReportInterval = "my int"
	assert.Equal(t, config, checkAgentEnvs(config))
}

func TestCheckServerEnvs(t *testing.T) {
	config, _ := NewServerEnvConfig()
	config.Address = "some addr"
	config.Key = "some key"
	config.Database = "my db"
	config.Restore = "true"
	config.StoreFile = "my file"
	config.StoreInterval = "some interval"
	assert.Equal(t, config, checkServerEnvs(config))
}
