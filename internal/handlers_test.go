package metrics

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler(t *testing.T) {
	// определяем структуру теста
	type want struct {
		code        int
		response    string
		contentType string
	}
	type arguments struct {
		metricType  string
		metricName  string
		metricValue string
	}
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name string
		want want
		args arguments
	}{
		// определяем все тесты
		{
			name: "Positive test Gauge 1",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
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
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
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
				code:        500,
				response:    "Wrong Gauge value\n",
				contentType: "text/plain; charset=utf-8",
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
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
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
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
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
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/plain",
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
				code:        500,
				response:    "Wrong Counter value\n",
				contentType: "text/plain; charset=utf-8",
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
				code:        500,
				response:    "Wrong Counter value\n",
				contentType: "text/plain; charset=utf-8",
			},
			args: arguments{
				metricType:  "counter",
				metricName:  "newValue2",
				metricValue: "dfdf",
			},
		},
	}
	metricsHandler := MetricsHandler{
		Metrics: NewMetrics(),
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			urlPath := fmt.Sprintf("/update/%s/%s/%s", tt.args.metricType, tt.args.metricName, tt.args.metricValue)
			request := httptest.NewRequest(http.MethodPost, urlPath, nil)
			request.Header.Set("Content-Type", "text/plain")

			// создаём новый Recorder
			w := httptest.NewRecorder()

			metricsHandler.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
