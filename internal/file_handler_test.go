package metrics

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileReaderWriter(t *testing.T) {
	filename := "test.json"
	testCountersMap := map[string]Counter{"PollCount": 1}
	testGaugeMap := map[string]Gauge{"RandomValue": 222.22, "Alloc": 11.11, "Frees": 33.3}
	testMetric := Metrics{
		CounterMetrics: testCountersMap,
		GaugeMetrics:   testGaugeMap,
	}
	defer func() {
		err := os.Remove(filename)
		if err != nil {
			panic(err)
		}
	}()
	testWriter, err := NewFileWriter(filename)
	defer func() {
		err := testWriter.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	testReader, err := NewFileReader(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := testReader.Close()
		if err != nil {
			panic(err)
		}
	}()
	assert.NotPanics(t, func() { testWriter.WriteJson(&testMetric) })
	jsonContent, err := testReader.ReadJson()
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonContent)
	assert.Equal(t, jsonContent, &testMetric)
}
