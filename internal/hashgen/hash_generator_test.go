package hashgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashGenerator(t *testing.T) {
	type want struct {
		hash []byte
	}
	type arguments struct {
		metricType string
		id         string
		value      interface{}
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Positive test HashGenerator 1",
			want: want{
				hash: []byte{0x30, 0x35, 0x63, 0x33, 0x63, 0x32, 0x65, 0x32, 0x64, 0x34, 0x30, 0x32, 0x30, 0x36, 0x35, 0x30, 0x61, 0x61, 0x37, 0x30, 0x61, 0x64, 0x62, 0x30, 0x35, 0x30, 0x39, 0x30, 0x34, 0x30, 0x38, 0x32, 0x65, 0x66, 0x38, 0x65, 0x65, 0x66, 0x62, 0x38, 0x31, 0x37, 0x30, 0x31, 0x62, 0x31, 0x31, 0x31, 0x64, 0x32, 0x66, 0x36, 0x36, 0x37, 0x66, 0x62, 0x30, 0x65, 0x36, 0x65, 0x37, 0x34, 0x35, 0x62},
			},
			args: arguments{
				metricType: "gauge",
				id:         "Alloc",
				value:      4.33,
			},
		},
		{
			name: "Positive test HashGenerator 2",
			want: want{
				hash: []byte{0x63, 0x63, 0x33, 0x39, 0x37, 0x30, 0x66, 0x33, 0x66, 0x63, 0x32, 0x38, 0x35, 0x64, 0x61, 0x39, 0x30, 0x61, 0x62, 0x66, 0x66, 0x39, 0x35, 0x31, 0x31, 0x33, 0x61, 0x32, 0x35, 0x62, 0x36, 0x31, 0x33, 0x62, 0x34, 0x39, 0x30, 0x32, 0x61, 0x32, 0x36, 0x66, 0x36, 0x63, 0x61, 0x37, 0x37, 0x31, 0x35, 0x30, 0x63, 0x39, 0x32, 0x32, 0x35, 0x66, 0x32, 0x35, 0x66, 0x33, 0x65, 0x39, 0x37, 0x65},
			},
			args: arguments{
				metricType: "counter",
				id:         "PollCount",
				value:      444,
			},
		},
	}
	generator := HashGenerator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generatedHash := generator.GenerateHash(tt.args.metricType, tt.args.id, tt.args.value)
			assert.Equal(t, []byte(generatedHash), tt.want.hash)
		})
	}
}