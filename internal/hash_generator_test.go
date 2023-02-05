package metrics

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
				hash: []byte{0xa3, 0x9, 0x5f, 0xfe, 0xcc, 0x9a, 0xc6, 0x2d, 0xf1, 0x24, 0x86, 0x8, 0x85, 0xb5, 0xb2, 0x1d, 0xe, 0x9, 0x2, 0x8, 0x5a, 0x41, 0xa8, 0x6b, 0x80, 0x15, 0x15, 0xa3, 0x8f, 0xa2, 0x60, 0xbe},
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
				hash: []byte{0x8, 0x39, 0x12, 0x7a, 0x7c, 0x77, 0xf5, 0x5a, 0x4f, 0xac, 0xde, 0x49, 0xc0, 0xd, 0xb2, 0x6, 0x3, 0x29, 0x87, 0x1f, 0x5f, 0x54, 0xbc, 0xe0, 0xac, 0x85, 0x3e, 0xb4, 0x53, 0x10, 0xb4, 0x1d},
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
