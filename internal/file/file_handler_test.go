package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReaderWriter(t *testing.T) {
	filename := "test.json"
	testCountersMap := map[string]metrics.Counter{"PollCount": 1}
	testGaugeMap := map[string]metrics.Gauge{"RandomValue": 222.22, "Alloc": 11.11, "Frees": 33.3}
	testMetric := metrics.Metrics{
		CounterMetrics: testCountersMap,
		GaugeMetrics:   testGaugeMap,
	}
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}()
	testWriter, err := NewWriter(filename)
	if err != nil {
		panic(err)
	}
	defer func(writer Writer) {
		err := writer.Close()
		if err != nil {
			panic(err)
		}
	}(*testWriter)

	testReader, err := NewReader(filename)
	if err != nil {
		panic(err)
	}
	defer func(reader Reader) {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}(*testReader)

	assert.NotPanics(t, func() { testWriter.WriteJSON(&testMetric) })
	jsonContent, err := testReader.ReadJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonContent)
	assert.Equal(t, jsonContent, &testMetric)
}

func TestReaderWriterNegative(t *testing.T) {
	filename := ""

	_, err := NewWriter(filename)
	require.Error(t, err)

	_, err = NewReader(filename)
	require.Error(t, err)
}
