package metrics

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func TestPOSTMetricsHandler(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	type arguments struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Positive test Gauge 1",
			want: want{
				code:     200,
				response: `{"status":"ok"}`,
			},
			args: arguments{
				metricType:  "gauge",
				metricName:  "Alloc",
				metricValue: "4.0",
			},
		},
		{
			name: "Positive test Gauge 2",
			want: want{
				code:     200,
				response: `{"status":"ok"}`,
			},
			args: arguments{
				metricType:  "gauge",
				metricName:  "testGauge",
				metricValue: "5.55555",
			},
		},
		{
			name: "Negative test Gauge 1",
			want: want{
				code:     400,
				response: "Wrong Gauge value\n",
			},
			args: arguments{
				metricType:  "gauge",
				metricName:  "testGauge",
				metricValue: "dsfsd",
			},
		},
		{
			name: "Positive test Gauge 3",
			want: want{
				code:     200,
				response: `{"status":"ok"}`,
			},
			args: arguments{
				metricType:  "gauge",
				metricName:  "11111",
				metricValue: "333",
			},
		},
		{
			name: "Positive test Counter 1",
			want: want{
				code:     200,
				response: `{"status":"ok"}`,
			},
			args: arguments{
				metricType:  "counter",
				metricName:  "testCounter",
				metricValue: "3",
			},
		},
		{
			name: "Positive test Counter 2",
			want: want{
				code:     200,
				response: `{"status":"ok"}`,
			},
			args: arguments{
				metricType:  "counter",
				metricName:  "myValue",
				metricValue: "100",
			},
		},
		{
			name: "Negative test Counter 1",
			want: want{
				code:     400,
				response: "Wrong Counter value\n",
			},
			args: arguments{
				metricType:  "counter",
				metricName:  "newValue",
				metricValue: "44444.444",
			},
		},
		{
			name: "Negative test Counter 2",
			want: want{
				code:     400,
				response: "Wrong Counter value\n",
			},
			args: arguments{
				metricType:  "counter",
				metricName:  "newValue2",
				metricValue: "dfdf",
			},
		},
	}
	var collector = NewCollector()
	metricsHandler := NewHandler(collector)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/update/%s/%s/%s", tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			statusCode, body := testRequest(t, ts, "POST", urlPath)
			assert.Equal(t, tt.want.code, statusCode)
			assert.Equal(t, tt.want.response, body)
		})
	}
}

func TestGETMetricsHandler(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	type arguments struct {
		metricType string
		metricName string
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Positive test Gauge 1",
			want: want{
				code:     200,
				response: "0",
			},
			args: arguments{
				metricType: "gauge",
				metricName: "Alloc",
			},
		},
		{
			name: "Negative test Gauge 1",
			want: want{
				code:     404,
				response: "Metric not found\n",
			},
			args: arguments{
				metricType: "gauge",
				metricName: "testGauge44",
			},
		},
		{
			name: "Positive test Gauge 2",
			want: want{
				code:     200,
				response: "0",
			},
			args: arguments{
				metricType: "gauge",
				metricName: "Frees",
			},
		},
		{
			name: "Positive test Counter 1",
			want: want{
				code:     200,
				response: "0",
			},
			args: arguments{
				metricType: "counter",
				metricName: "PollCount",
			},
		},
		{
			name: "Negative test Counter 1",
			want: want{
				code:     404,
				response: "Metric not found\n",
			},
			args: arguments{
				metricType: "counter",
				metricName: "mymetric",
			},
		},
	}
	var collector = NewCollector()
	metricsHandler := NewHandler(collector)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/value/%s/%s", tt.args.metricType, tt.args.metricName)
			statusCode, body := testRequest(t, ts, "GET", urlPath)
			assert.Equal(t, tt.want.code, statusCode)
			assert.Equal(t, tt.want.response, body)
		})
	}
}
