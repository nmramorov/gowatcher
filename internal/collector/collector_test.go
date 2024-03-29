package collector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/errors"
)

func TestNewCollector(t *testing.T) {
	newCollector := NewCollector()
	assert.EqualValues(t, newCollector.Updates, 0)
}

func TestUpdateMetrics(t *testing.T) {
	newCollector := NewCollector()
	newCollector.UpdateMetrics()
	assert.Equal(t, newCollector.Updates, 1)
	assert.Equal(t, newCollector.Metrics.CounterMetrics["PollCount"], metrics.Counter(1))

	newCollector.UpdateMetrics()
	newCollector.UpdateMetrics()
	assert.Equal(t, newCollector.Metrics.CounterMetrics["PollCount"], metrics.Counter(3))
}

func TestGetMetrics(t *testing.T) {
	newCollector := NewCollector()
	newCollector.UpdateMetrics()
	_metrics := newCollector.GetMetrics()
	assert.Equal(t, _metrics.CounterMetrics["PollCount"], metrics.Counter(1))
}

func TestGetMetricSuccess(t *testing.T) {
	newCollector := NewCollector()
	result, err := newCollector.GetMetric("PollCount")
	if err != nil {
		errLog := fmt.Errorf("wrong result in GetMetric %w", err)
		fmt.Println(errLog)
	}
	assert.Equal(t, result, metrics.Counter(0))
}

func TestGetMetricError(t *testing.T) {
	newCollector := NewCollector()
	_, err := newCollector.GetMetric("SampleKey")
	assert.Equal(t, err, errors.ErrorMetricNotFound)
}

func TestStringFromCounter(t *testing.T) {
	newCollector := NewCollector()
	pollCount, err := newCollector.GetMetric("PollCount")
	if err != nil {
		panic(1)
	}
	strPollCount := fmt.Sprint(pollCount)
	assert.Equal(t, strPollCount, "0")
}

func TestStringFromGauge(t *testing.T) {
	newCollector := NewCollector()
	alloc, err := newCollector.GetMetric("Alloc")
	if err != nil {
		panic(1)
	}
	strAlloc := fmt.Sprint(alloc)
	assert.GreaterOrEqual(t, strAlloc, "0")
}

func TestUpdateBatch(t *testing.T) {
	newCollector := NewCollector()
	alloc := float64(2.2)
	pollCount := int64(44444)
	myVal := float64(3333.3333)
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

func TestGetMetric(t *testing.T) {
	newCollector := NewCollector()
	alloc := float64(2.2)
	pollCount := int64(44444)
	myVal := float64(3333.3333)
	newCollector.UpdateMetrics()
	toUpdate := []*metrics.JSONMetrics{
		{
			ID:    "Alloc",
			MType: "gauge",
			Value: &alloc,
			Delta: nil,
		},
		{
			ID:    "MyMetric",
			MType: "gauge",
			Value: &myVal,
			Delta: nil,
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Value: nil,
			Delta: &pollCount,
		},
	}
	newCollector.UpdateBatch(toUpdate)
	res, _ := newCollector.GetMetricJSON(toUpdate[0])
	assert.Equal(t, alloc, *res.Value)
	res, _ = newCollector.GetMetricJSON(toUpdate[1])
	assert.Equal(t, myVal, *res.Value)
	res, _ = newCollector.GetMetricJSON(toUpdate[2])
	assert.Equal(t, pollCount+1, *res.Delta)
}

func TestNewCollecorFromSavedFile(t *testing.T) {
	assert.NotPanics(t, func() {
		saved := metrics.NewMetrics()
		_ = NewCollectorFromSavedFile(saved)
	})
}

func TestCollectorString(t *testing.T) {
	c := NewCollector()
	testFloat64, _ := c.String(float64(22.0))
	assert.Equal(t, "22", testFloat64)
	testInt64, _ := c.String(int64(22))
	assert.Equal(t, "22", testInt64)
}

func TestCollectorUpdateExtraMetrics(t *testing.T) {
	c := NewCollector()
	c.UpdateExtraMetrics()
	assert.GreaterOrEqual(t, c.GetMetrics().GaugeMetrics["TotalMemory"], metrics.Gauge(0.0))
	assert.GreaterOrEqual(t, c.GetMetrics().GaugeMetrics["FreeMemory"], metrics.Gauge(0.0))
	assert.GreaterOrEqual(t, c.GetMetrics().GaugeMetrics["CPUutilization1"], metrics.Gauge(0.0))
}
