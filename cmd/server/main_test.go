package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"internal/metrics"
)

func TestServer(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	type arguments struct {
		url string
	}
	tests := []struct {
		name string
		want want
		args arguments
	}{
		{
			name: "Test Positive Gauge",
			want: want{
				code:     200,
				response: `{"status":"ok"}`,
			},
			args: arguments{
				url: "http://localhost:8080/update/gauge/newMetric/100.11",
			},
		},
		{
			name: "Test Negative Gauge",
			want: want{
				code:     400,
				response: "Wrong Gauge value\n",
			},
			args: arguments{
				url: "http://localhost:8080/update/gauge/newMetric/none",
			},
		},
		{
			name: "Test Positive Counter",
			want: want{
				code:     200,
				response: `{"status":"ok"}`,
			},
			args: arguments{
				url: "http://localhost:8080/update/counter/newMetric/100",
			},
		},
		{
			name: "Test Negative Counter",
			want: want{
				code:     400,
				response: "Wrong Counter value\n",
			},
			args: arguments{
				url: "http://localhost:8080/update/counter/newMetric/none",
			},
		},
		{
			name: "Test wrong method 1",
			want: want{
				code:     404,
				response: "404 page not found\n",
			},
			args: arguments{
				url: "http://localhost:8080/do_something_new",
			},
		},
		{
			name: "Test wrong method 2",
			want: want{
				code:     404,
				response: "404 page not found\n",
			},
			args: arguments{
				url: "http://localhost:8080/updater/gauge/newMetric/none",
			},
		},
		{
			name: "Test wrong path 1",
			want: want{
				code:     404,
				response: "404 page not found\n",
			},
			args: arguments{
				url: "http://localhost:8080/update/gauge/",
			},
		},
		{
			name: "Test wrong path 1",
			want: want{
				code:     404,
				response: "404 page not found\n",
			},
			args: arguments{
				url: "http://localhost:8080/update/counter/dsfsdf/dsfdsff/dsfsd/fsd",
			},
		},
	}
	MOCK_CURSOR := metrics.NewCursor("", "pgx")
	metricsHandler := metrics.NewHandler("", MOCK_CURSOR)

	ts := httptest.NewServer(metricsHandler)

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.args.url, nil)
			request.Header.Set("Content-Type", "text/plain")

			w := httptest.NewRecorder()

			metricsHandler.ServeHTTP(w, request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
		})
	}
}
