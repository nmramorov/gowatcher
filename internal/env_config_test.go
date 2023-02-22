package metrics

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
	assert.Equal(t, testConfig.RateLimit, "")
}
