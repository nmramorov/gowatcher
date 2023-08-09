package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"testing"

	"github.com/nmramorov/gowatcher/internal/collector"
	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/db"
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
		{
			name: "Negative test Counter 3",
			want: want{
				code:     405,
				response: "",
			},
			args: arguments{
				metricType:  "counter",
				metricName:  "newValue2",
				metricValue: "dfdf",
			},
		},
	}
	ctx := context.Background()
	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("", "", "", MOCKCURSOR)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/update/%s/%s/%s", tt.args.metricType,
				tt.args.metricName, tt.args.metricValue)
			if tt.name == "Negative test Counter 3" {
				statusCode, body := testRequest(t, ts, "GET", urlPath)
				assert.Equal(t, tt.want.code, statusCode)
				assert.Equal(t, tt.want.response, body)
				return
			}
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
	ctx := context.Background()
	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("", "", "", MOCKCURSOR)

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
	ctx := context.Background()

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("", "", "", MOCKCURSOR)
	metricsHandler.Collector.UpdateMetrics()

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
	ctx := context.Background()

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("", "", "", MOCKCURSOR)

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
	ctx := context.Background()

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("", "", "", MOCKCURSOR)

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
			metricsHandler.Collector.UpdateMetrics()
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
	ctx := context.Background()

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("fdsfds", "", "", MOCKCURSOR)

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

func TestNewHandlerFromSavedData(t *testing.T) {
	ctx := context.Background()
	col := collector.NewCollector()
	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	assert.NotPanics(t, func() {
		NewHandlerFromSavedData(col.GetMetrics(),
			"sss", "", "", MOCKCURSOR)
	})
}

func TestPOSTMetricsHandlerJsonBatch(t *testing.T) {
	GaugeVal := 44.4
	var CountVal int64 = 3

	type want struct {
		code int
	}
	tests := []struct {
		name string
		want want
		args []m.JSONMetrics
	}{
		{
			name: "Negative test Gauge 1",
			want: want{
				code: 400,
			},
			args: []m.JSONMetrics{
				{
					ID:    "GaugeMetric",
					MType: "gauge",
					Value: &GaugeVal,
				},
			},
		},
		{
			name: "Negative test Counter 1",
			want: want{
				code: 400,
			},
			args: []m.JSONMetrics{
				{
					ID:    "CounterMetric",
					MType: "counter",
					Delta: &CountVal,
				},
			},
		},
	}
	ctx := context.Background()

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("", "", "", MOCKCURSOR)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/update/"
			statusCode, body := testRequestJSON(t, ts, "POST", urlPath, tt.args)
			assert.Equal(t, tt.want.code, statusCode)
			assert.NotNil(t, body)
		})
	}
}

func TestPOSTUpdateJsonBatch(t *testing.T) {
	GaugeVal := 44.4
	var CountVal int64 = 3

	type want struct {
		code int
	}
	tests := []struct {
		name string
		want want
		args []m.JSONMetrics
	}{
		{
			name: "Positive test Gauge 1",
			want: want{
				code: 200,
			},
			args: []m.JSONMetrics{
				{
					ID:    "GaugeMetric",
					MType: "gauge",
					Value: &GaugeVal,
				},
			},
		},
		{
			name: "Positive test Counter 1",
			want: want{
				code: 200,
			},
			args: []m.JSONMetrics{
				{
					ID:    "CounterMetric",
					MType: "counter",
					Delta: &CountVal,
				},
			},
		},
	}
	ctx := context.Background()

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("", "", "", MOCKCURSOR)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/updates/"
			statusCode, body := testRequestJSON(t, ts, "POST", urlPath, tt.args)
			assert.Equal(t, tt.want.code, statusCode)
			assert.NotNil(t, body)
		})
	}
}

func TestDecodeMsg(t *testing.T) {
	GaugeVal := 44.4
	var CountVal int64 = 3
	ctx := context.Background()
	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")

	type want struct {
		code int
	}
	tests := []struct {
		name    string
		want    want
		handler *Handler
		args    *m.JSONMetrics
	}{
		{
			name: "Negative test Gauge 1",
			want: want{
				code: 400,
			},
			handler: NewHandler("dfd", "./key.pem", "", MOCKCURSOR),
			args: &m.JSONMetrics{
				ID:    "GaugeMetric",
				MType: "gauge",
				Value: &GaugeVal,
			},
		},
		{
			name:    "Negative test Counter 1",
			handler: NewHandler("dsfdsf", "./key.pem", "", MOCKCURSOR),
			want: want{
				code: 400,
			},
			args: &m.JSONMetrics{
				ID:    "CounterMetric",
				MType: "counter",
				Delta: &CountVal,
			},
		},
		{
			name:    "Positive test Counter 1",
			handler: NewHandler("dfdsfsdf", "", "", MOCKCURSOR),
			want: want{
				code: 200,
			},
			args: &m.JSONMetrics{
				ID:    "CounterMetric",
				MType: "counter",
				Delta: &CountVal,
			},
		},
		{
			name:    "Negative test Counter 2",
			handler: NewHandler("dsfsdf", "./.", "", MOCKCURSOR),
			want: want{
				code: 400,
			},
			args: &m.JSONMetrics{
				ID:    "CounterMetric",
				MType: "counter",
				Delta: &CountVal,
			},
		},
		{
			name:    "Negative test Counter 3",
			handler: NewHandler("sdfdsfds", "./key.pem", "", MOCKCURSOR),
			want: want{
				code: 400,
			},
			args: &m.JSONMetrics{
				ID:    "CounterMetric",
				MType: "fdsf",
				Delta: &CountVal,
			},
		},
	}

	for _, tt := range tests {
		ts := httptest.NewServer(tt.handler)

		defer ts.Close()
		t.Run(tt.name, func(t *testing.T) {
			urlPath := "/update/"
			statusCode, body := testRequestJSON(t, ts, "POST", urlPath, tt.args)
			assert.Equal(t, tt.want.code, statusCode)
			assert.NotNil(t, body)
		})
	}
}

func TestHandler_ValidateIP(t *testing.T) {
	ctx := context.Background()

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandlerWithoutSubnet := NewHandler("", "", "", MOCKCURSOR)
	metricsHandlerWithWrongCIDR := NewHandler("", "", "sdfdsf", MOCKCURSOR)
	metricsHandlerWithProperSubnet := NewHandler("", "", "192.168.1.0/24", MOCKCURSOR)
	metricsHandlerWithWrongSubnet := NewHandler("", "", "192.168.0.1/24", MOCKCURSOR)

	type want struct {
		code int
	}
	tests := []struct {
		name string
		want want
		args *Handler
	}{
		{
			name: "Positive test no subnet provided",
			want: want{
				code: 200,
			},
			args: metricsHandlerWithoutSubnet,
		},
		{
			name: "Positive test proper subnet provided",
			want: want{
				code: 200,
			},
			args: metricsHandlerWithProperSubnet,
		},
		{
			name: "Negative test wrong subnet",
			want: want{
				code: 403,
			},
			args: metricsHandlerWithWrongSubnet,
		},
		{
			name: "Negative test wrong CIDR",
			want: want{
				code: 500,
			},
			args: metricsHandlerWithWrongCIDR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.args)

			defer ts.Close()
			urlPath := ts.URL + "/"
			req, err := http.NewRequest("GET", urlPath, nil)
			req.Header.Add("X-Real-IP", "192.168.1.11:8080")
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}

func TestCheckHash(t *testing.T) {
	ctx := context.Background()
	var delta int64 = 3
	var value float64 = 0.0

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("very secret key", "", "", MOCKCURSOR)
	err := metricsHandler.CheckHash(&m.JSONMetrics{
		ID:    "MyMetric",
		MType: "counter",
		Delta: &delta,
	})
	require.NoError(t, err)
	err = metricsHandler.CheckHash(&m.JSONMetrics{
		ID:    "MyMetric",
		MType: "gauge",
		Value: &value,
	})
	require.NoError(t, err)
}

func TestInitDB(t *testing.T) {
	ctx := context.Background()

	MOCKCURSOR, _ := db.NewCursor(ctx, "", "pgx")
	metricsHandler := NewHandler("very secret key", "", "", MOCKCURSOR)
	err := metricsHandler.InitDB(ctx)
	require.Error(t, err)
}
