package metrics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCollector(t *testing.T) {
	var newCollector = NewCollector()
	assert.EqualValues(t, newCollector.updates, 0)
}

func TestUpdateMetrics(t *testing.T) {
	var newCollector = NewCollector()
	newCollector.UpdateMetrics()
	assert.Equal(t, newCollector.updates, 1)
	assert.Equal(t, newCollector.metrics.CountMetrics["PollCount"], counter(1))

	newCollector.UpdateMetrics()
	newCollector.UpdateMetrics()
	assert.Equal(t, newCollector.metrics.CountMetrics["PollCount"], counter(3))
}

func TestGetMetrics(t *testing.T) {
	var newCollector = NewCollector()
	newCollector.UpdateMetrics()
	metrics := newCollector.GetMetrics()
	assert.Equal(t, metrics.CountMetrics["PollCount"], counter(1))
}

func TestGetMetricSuccess(t *testing.T) {
	var newCollector = NewCollector()
	result, err := newCollector.GetMetric("PollCount")
	if err != nil {
		errLog := fmt.Errorf("wrong result in GetMetric %w", err)
		fmt.Println(errLog)
	}
	assert.Equal(t, result, counter(0))
}

func TestGetMetricError(t *testing.T) {
	var newCollector = NewCollector()
	_, err := newCollector.GetMetric("SampleKey")
	assert.Equal(t, err, ErrorMetricNotFound)
}
