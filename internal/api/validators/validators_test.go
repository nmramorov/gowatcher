package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	col "github.com/nmramorov/gowatcher/internal/collector"
)

type test struct {
	name      string
	inputType string
	input     string
	result    bool
}

var tests []test = []test{
	{
		name:      "Positive Validation 1",
		inputType: "gauge",
		input:     "Alloc",
		result:    true,
	},
	{
		name:      "Positive Validation 2",
		inputType: "counter",
		input:     "PollCount",
		result:    true,
	},
	{
		name:      "Negative Validation 1",
		inputType: "counter",
		input:     "MyCounterMetric",
		result:    false,
	},
	{
		name:      "Negative Validation 2",
		inputType: "gauge",
		input:     "MyGaugeMetric",
		result:    false,
	},
	{
		name:      "Negative Validation 3",
		inputType: "sample",
		input:     "SampleValue",
		result:    false,
	},
}

func TestValidateMetricType(t *testing.T) {
	newCollector := col.NewCollector()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ValidateMetric(tt.inputType, tt.input, newCollector), tt.result)
		})
	}
}
