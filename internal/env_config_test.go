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
	assert.Equal(t, testConfig.Address, ``)
	assert.Equal(t, testConfig.ReportInterval, ``)
	assert.Equal(t, testConfig.PollInterval, ``)
	assert.Equal(t, testConfig.StoreInterval, ``)
	assert.Equal(t, testConfig.StoreFile, "")
	assert.Equal(t, testConfig.Restore, false)
}

func TestEnvConfigIntervalConvertion(t *testing.T) {
	testConfig, err := NewConfig()
	if err != nil {
		panic(err)
	}
	testConfig.PollInterval = "2s"
	testConfig.ReportInterval = "10m"
	reportInterval, err := testConfig.GetNumericInterval("ReportInterval")
	fmt.Println(reportInterval)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, reportInterval, int64(600))
	pollInterval, err := testConfig.GetNumericInterval("PollInterval")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, pollInterval, int64(2))
}
