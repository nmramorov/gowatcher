package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvConfig(t *testing.T) {
	testConfig, _ := NewConfig()
	assert.Equal(t, testConfig.Address, `127.0.0.1:8080`)
	assert.Equal(t, testConfig.ReportInterval, `10s`)
	assert.Equal(t, testConfig.PollInterval, `2s`)
	assert.Equal(t, testConfig.StoreInterval, `300s`)
	assert.Equal(t, testConfig.StoreFile, "/tmp/devops-metrics-db.json")
	assert.Equal(t, testConfig.Restore, true)
}

func TestEnvConfigIntervalConvertion(t *testing.T) {
	testConfig, _ := NewConfig()
	testConfig.PollInterval = "2s"
	testConfig.ReportInterval = "10m"
	reportInterval := testConfig.GetNumericInterval("ReportInterval")
	assert.Equal(t, reportInterval, int64(600))
	pollInterval := testConfig.GetNumericInterval("PollInterval")
	assert.Equal(t, pollInterval, int64(2))
}
