package metrics

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testMetricStr string = `
id:1,
type:counter,
delta:4\n
id:1,
type:counter,
delta:4\n
`

func TestFileReaderWriter(t *testing.T) {
	filename := "test.json"
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
	var testDelta int64 = 4
	testMetric := JSONMetrics{
		ID:    "1",
		MType: "counter",
		Delta: &testDelta,
	}
	assert.NotPanics(t, func() { testWriter.WriteJson(&testMetric) })
	jsonContent, err := testReader.ReadJson()
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonContent)
	assert.Equal(t, "\nid:1,\ntype:counter,\ndelta:4\\n\nid:1,\ntype:counter,\ndelta:4\\n\n", testMetricStr)
}
