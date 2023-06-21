package handlers

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

	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/db"
	// "github.com/nmramorov/gowatcher/internal/log"
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

func testRequestJSON(t *testing.T, ts *httptest.Server, method, path string, payload interface{}) (int, []byte) {
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
	MOCKCURSOR, _ := db.NewCursor("", "pgx")
	metricsHandler := NewHandler("", MOCKCURSOR)

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
	MOCKCURSOR, _ := db.NewCursor("", "pgx")
	metricsHandler := NewHandler("", MOCKCURSOR)

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
	MOCKCURSOR, _ := db.NewCursor("", "pgx")
	metricsHandler := NewHandler("", MOCKCURSOR)
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
		response m.JSONMetrics
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
				response: m.JSONMetrics{
					ID:    "GaugeMetric",
					MType: "gauge",
					Value: &GaugeVal,
					Hash:  "b80631b192b6e327725bf20e38fb0ca59cf84026515fbc8b3e2a9727ace1e313",
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
				response: m.JSONMetrics{
					ID:    "CounterMetric",
					MType: "counter",
					Delta: &CountVal,
					Hash:  "34e44bc45730d6ead6c5959c68a8d591f932afac6522a71df1bea414deb21fdd",
				},
			},
			args: arguments{
				ID:    "CounterMetric",
				MType: "counter",
				Delta: &CountVal,
			},
		},
	}

	MOCKCURSOR, _ := db.NewCursor("", "pgx")
	metricsHandler := NewHandler("", MOCKCURSOR)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/update/"
			statusCode, body := testRequestJSON(t, ts, "POST", urlPath, tt.args)
			assert.Equal(t, tt.want.code, statusCode)
			result := m.JSONMetrics{}
			if err := json.Unmarshal(body, &result); err != nil {
				panic("Test error: error with unmarshalling JSON POST /update method")
			}
			assert.Equal(t, tt.want.response, result)
		})
	}
}

// ToDO: Fix Example!
// func Example() {
// 	GaugeVal := 44.4

// 	type arguments struct {
// 		ID    string   `json:"id"`
// 		MType string   `json:"type"`
// 		Delta *int64   `json:"delta,omitempty"`
// 		Value *float64 `json:"value,omitempty"`
// 	}

// 	MOCKCURSOR, _ := db.NewCursor("", "pgx")
// 	metricsHandler := NewHandler("", MOCKCURSOR)

// 	ts := httptest.NewServer(metricsHandler)

// 	defer ts.Close()

// 	payload := arguments{
// 		ID:    "GaugeMetric",
// 		MType: "gauge",
// 		Value: &GaugeVal,
// 	}

// 	buf := bytes.NewBuffer([]byte{})
// 	encoder := json.NewEncoder(buf)
// 	encoder.Encode(payload)
// 	req, err := http.NewRequest("POST", ts.URL+"/update/", bytes.NewBuffer(buf.Bytes()))
// 	if err != nil {
// 		log.ErrorLog.Printf("Error occured in example: %e", err)
// 	}
// 	req.Header.Add("Content-Type", "application/json")

// 	resp, err := http.DefaultClient.Do(req)

// 	respBody, err := io.ReadAll(resp.Body)

// 	defer resp.Body.Close()

// 	result := m.JSONMetrics{}
// 	if err := json.Unmarshal(respBody, &result); err != nil {
// 		panic("Test error: error with unmarshalling JSON POST /update method")
// 	}
// 	fmt.Println(result)

// 	// Output:
// 	// JSONMetrics{
// 	//	ID:    "GaugeMetric",
// 	//	MType: "gauge",
// 	//	Value: &GaugeVal,
// 	//	Hash:  "b80631b192b6e327725bf20e38fb0ca59cf84026515fbc8b3e2a9727ace1e313",
// 	// },
// }

func TestPOSTValueMetricsHandlerJson(t *testing.T) {
	GaugeVal := 0.0
	var CountVal int64 = 0
	var PollCount1 int64 = 2
	var PollCount2 int64 = 3

	type want struct {
		code     int
		response m.JSONMetrics
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
				response: m.JSONMetrics{
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
				response: m.JSONMetrics{
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
				response: m.JSONMetrics{
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
				response: m.JSONMetrics{
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

	MOCKCURSOR, _ := db.NewCursor("", "pgx")
	metricsHandler := NewHandler("", MOCKCURSOR)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/value/"
			statusCode, body := testRequestJSON(t, ts, "POST", urlPath, tt.args)
			assert.Equal(t, tt.want.code, statusCode)
			result := m.JSONMetrics{}
			if err := json.Unmarshal(body, &result); err != nil {
				panic("Test error: error with unmarshalling JSON POST /value method")
			}
			assert.Equal(t, tt.want.response, result)
			metricsHandler.collector.UpdateMetrics()
		})
	}
}

func TestPing(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Negative test Ping",
			want: want{
				code:     500,
				response: ``,
			},
		},
	}
	MOCKCURSOR, _ := db.NewCursor("", "pgx")
	metricsHandler := NewHandler("", MOCKCURSOR)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/ping"
			statusCode, resp := testRequest(t, ts, "GET", urlPath)
			fmt.Println(resp)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}
