package collector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/errors"
)

func TestNewCollector(t *testing.T) {
	var newCollector = NewCollector()
	assert.EqualValues(t, newCollector.Updates, 0)
}

func TestUpdateMetrics(t *testing.T) {
	var newCollector = NewCollector()
	newCollector.UpdateMetrics()
	assert.Equal(t, newCollector.Updates, 1)
	assert.Equal(t, newCollector.Metrics.CounterMetrics["PollCount"], metrics.Counter(1))

	newCollector.UpdateMetrics()
	newCollector.UpdateMetrics()
	assert.Equal(t, newCollector.Metrics.CounterMetrics["PollCount"], metrics.Counter(3))
}

func TestGetMetrics(t *testing.T) {
	var newCollector = NewCollector()
	newCollector.UpdateMetrics()
	_metrics := newCollector.GetMetrics()
	assert.Equal(t, _metrics.CounterMetrics["PollCount"], metrics.Counter(1))
}

func TestGetMetricSuccess(t *testing.T) {
	var newCollector = NewCollector()
	result, err := newCollector.GetMetric("PollCount")
	if err != nil {
		errLog := fmt.Errorf("wrong result in GetMetric %w", err)
		fmt.Println(errLog)
	}
	assert.Equal(t, result, metrics.Counter(0))
}

func TestGetMetricError(t *testing.T) {
	var newCollector = NewCollector()
	_, err := newCollector.GetMetric("SampleKey")
	assert.Equal(t, err, errors.ErrorMetricNotFound)
}

func TestStringFromCounter(t *testing.T) {
	var newCollector = NewCollector()
	pollCount, err := newCollector.GetMetric("PollCount")
	if err != nil {
		panic(1)
	}
	strPollCount := fmt.Sprint(pollCount)
	if err != nil {
		panic(1)
	}
	assert.Equal(t, strPollCount, "0")
}

func TestStringFromGauge(t *testing.T) {
	var newCollector = NewCollector()
	alloc, err := newCollector.GetMetric("Alloc")
	if err != nil {
		panic(1)
	}
	strAlloc := fmt.Sprint(alloc)
	assert.GreaterOrEqual(t, strAlloc, "0")
}

func TestUpdateBatch(t *testing.T) {
	var newCollector = NewCollector()
	var alloc float64 = 2.2
	var pollCount int64 = 44444
	var myVal float64 = 3333.3333
	newCollector.UpdateMetrics()
	toUpdate := []*metrics.JSONMetrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &alloc,
			Delta: nil,
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Value: nil,
			Delta: &pollCount,
		},
		{
			ID:    "MyMetric",
			MType: "gauge",
			Value: &myVal,
			Delta: nil,
		},
	}
	assert.NotPanics(t, func() { newCollector.UpdateBatch(toUpdate) })
	assert.Equal(t, alloc, float64(newCollector.Metrics.GaugeMetrics["Alloc"]))
	assert.Equal(t, pollCount+1, int64(newCollector.Metrics.CounterMetrics["PollCount"]))
	assert.Equal(t, myVal, float64(newCollector.Metrics.GaugeMetrics["MyMetric"]))
}
