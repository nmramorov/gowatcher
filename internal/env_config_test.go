package metrics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvConfig(t *testing.T) {
	testConfig, err := NewConfig()
	if err != nil {
		panic(err)
	}
	assert.Equal(t, testConfig.Address, `127.0.0.1:8080`)
	assert.Equal(t, testConfig.ReportInterval, `10s`)
	assert.Equal(t, testConfig.PollInterval, `2s`)
	assert.Equal(t, testConfig.StoreInterval, `300s`)
	assert.Equal(t, testConfig.StoreFile, "/tmp/devops-metrics-db.json")
	assert.Equal(t, testConfig.Restore, true)
}

func TestEnvConfigIntervalConvertion(t *testing.T) {
	testConfig, err := NewConfig()
	if err != nil {
		panic(err)
	}
	reportInterval, err := testConfig.GetNumericInterval("ReportInterval")
	fmt.Println(reportInterval)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, reportInterval, int64(10))
	pollInterval, err := testConfig.GetNumericInterval("PollInterval")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, pollInterval, int64(2))
}
