package metrics

import (
	"bytes"
	"encoding/json"
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

func testRequestJson(t *testing.T, ts *httptest.Server, method, path string, payload interface{}) (int, []byte) {
	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	encoder.Encode(payload)
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(buf.Bytes()))
	req.Header.Add("Content-Type", "application/json")
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, respBody
}

func TestPOSTMetricsHandlerNoJson(t *testing.T) {
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

	metricsHandler := NewHandler("")

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
	metricsHandler := NewHandler("")

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

func TestHTML(t *testing.T) {
	metricsHandler := NewHandler("")
	metricsHandler.collector.UpdateMetrics()

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	statusCode, _ := testRequest(t, ts, "GET", "/")
	assert.Equal(t, 200, statusCode)
}

func TestPOSTMetricsHandlerJson(t *testing.T) {
	GaugeVal := 44.4
	var CountVal int64 = 3

	type want struct {
		code     int
		response JSONMetrics
	}
	type arguments struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Positive test Gauge 1",
			want: want{
				code: 200,
				response: JSONMetrics{
					ID:    "GaugeMetric",
					MType: "gauge",
					Value: &GaugeVal,
				},
			},
			args: arguments{
				ID:    "GaugeMetric",
				MType: "gauge",
				Value: &GaugeVal,
			},
		},
		{
			name: "Positive test Counter 1",
			want: want{
				code: 200,
				response: JSONMetrics{
					ID:    "CounterMetric",
					MType: "counter",
					Delta: &CountVal,
				},
			},
			args: arguments{
				ID:    "CounterMetric",
				MType: "counter",
				Delta: &CountVal,
			},
		},
	}

	metricsHandler := NewHandler("")

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/update/"
			statusCode, body := testRequestJson(t, ts, "POST", urlPath, tt.args)
			assert.Equal(t, tt.want.code, statusCode)
			result := JSONMetrics{}
			if err := json.Unmarshal(body, &result); err != nil {
				panic("Test error: error with unmarshalling JSON POST /update method")
			}
			assert.Equal(t, tt.want.response, result)
		})
	}
}

func TestPOSTValueMetricsHandlerJson(t *testing.T) {
	GaugeVal := 0.0
	var CountVal int64 = 0
	var PollCount1 int64 = 2
	var PollCount2 int64 = 3

	type want struct {
		code     int
		response JSONMetrics
	}
	type arguments struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Positive test Gauge 1",
			want: want{
				code: 200,
				response: JSONMetrics{
					ID:    "GaugeMetric",
					MType: "gauge",
					Value: &GaugeVal,
				},
			},
			args: arguments{
				ID:    "GaugeMetric",
				MType: "gauge",
			},
		},
		{
			name: "Positive test Counter 1",
			want: want{
				code: 200,
				response: JSONMetrics{
					ID:    "CounterMetric",
					MType: "counter",
					Delta: &CountVal,
				},
			},
			args: arguments{
				ID:    "CounterMetric",
				MType: "counter",
			},
		},
		{
			name: "Positive test Counter 2",
			want: want{
				code: 200,
				response: JSONMetrics{
					ID:    "PollCount",
					MType: "counter",
					Delta: &PollCount1,
				},
			},
			args: arguments{
				ID:    "PollCount",
				MType: "counter",
			},
		},
		{
			name: "Positive test Counter 3",
			want: want{
				code: 200,
				response: JSONMetrics{
					ID:    "PollCount",
					MType: "counter",
					Delta: &PollCount2,
				},
			},
			args: arguments{
				ID:    "PollCount",
				MType: "counter",
			},
		},
	}

	metricsHandler := NewHandler("")

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/value/"
			statusCode, body := testRequestJson(t, ts, "POST", urlPath, tt.args)
			assert.Equal(t, tt.want.code, statusCode)
			result := JSONMetrics{}
			if err := json.Unmarshal(body, &result); err != nil {
				panic("Test error: error with unmarshalling JSON POST /value method")
			}
			assert.Equal(t, tt.want.response, result)
			metricsHandler.collector.UpdateMetrics()
		})
	}
}
