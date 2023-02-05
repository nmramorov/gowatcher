package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"internal/metrics"
)

func CreateRequests(endpoint string, mtrcs *metrics.Metrics) []*http.Request {
	var requests []*http.Request
	for k, v := range mtrcs.GaugeMetrics {
		req, err := http.NewRequest(http.MethodPost, endpoint+fmt.Sprintf("/update/gauge/%s/%f", k, v), nil)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "text/plain")
		requests = append(requests, req)
	}
	for k, v := range mtrcs.CounterMetrics {
		req, err := http.NewRequest(http.MethodPost, endpoint+fmt.Sprintf("/update/counter/%s/%d", k, v), nil)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for counter with params: %s %d", k, v)
		}
		req.Header.Add("Content-Type", "text/plain")
		requests = append(requests, req)
	}
	return requests
}

func createBody(metricType, path, key, secretkey string, value interface{}) *bytes.Buffer {
	var hash string
	if secretkey != "" {
		generator := metrics.NewHashGenerator(secretkey)
		hash = generator.GenerateHash(metricType, key, value)
	}
	body := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(body)
	var toEncode metrics.JSONMetrics
	switch metricType {
	case "gauge":
		toEncode.MType = "gauge"
		toEncode.ID = key
		gaugeVal := value.(metrics.Gauge)
		val := (*float64)(&gaugeVal)
		toEncode.Value = val
		toEncode.Hash = hash
	case "counter":
		toEncode.MType = "counter"
		toEncode.ID = key
		counterVal := value.(metrics.Counter)
		val := (*int64)(&counterVal)
		toEncode.Delta = val
		toEncode.Hash = hash
	}
	toEncode.ID = key
	if path == "/value/" {
		toEncode.Delta = nil
		toEncode.Value = nil
	}
	encoder.Encode(&toEncode)
	return body
}

func createGaugeRequests(endpoint, path, key string, gaugeMetrics map[string]metrics.Gauge) []*http.Request {
	var requests []*http.Request
	for k, v := range gaugeMetrics {
		body := createBody("gauge", path, k, key, v)
		req, err := http.NewRequest(http.MethodPost, endpoint+path, body)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	return requests
}

func createCounterRequests(endpoint, path, key string, counterMetrics map[string]metrics.Counter) []*http.Request {
	var requests []*http.Request
	for k, v := range counterMetrics {
		body := createBody("counter", path, k, key, v)
		req, err := http.NewRequest(http.MethodPost, endpoint+path, body)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for counter with params: %s %d", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	return requests
}

func generateMetricsRequests(endpoint, path, key string, src *metrics.Metrics) []*http.Request {
	gaugeRequests := createGaugeRequests(endpoint, path, key, src.GaugeMetrics)
	counterRequests := createCounterRequests(endpoint, path, key, src.CounterMetrics)
	return append(gaugeRequests, counterRequests...)
}

func PushMetrics(client *http.Client, endpoint string, mtrcs *metrics.Metrics, key string) {
	defer func() {
		if p := recover(); p != nil {
			metrics.ErrorLog.Println(p)
		}
	}()
	requests := generateMetricsRequests(endpoint, "/update/", key, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			metrics.ErrorLog.Println(err)
		}
		defer resp.Body.Close()
	}
}

func GetMetricsValues(client *http.Client, endpoint, key string, mtrcs *metrics.Metrics) {
	defer func() {
		if p := recover(); p != nil {
			metrics.ErrorLog.Println(p)
		}
	}()
	requests := generateMetricsRequests(endpoint, "/value/", key, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			metrics.ErrorLog.Println(err)
		}
		defer resp.Body.Close()
	}
}

func main() {
	agentConfig := metrics.GetAgentConfig()
	fmt.Println(agentConfig)

	var collector = metrics.NewCollector()

	endpoint := "http://" + agentConfig.Address

	client := &http.Client{}
	metrics.InfoLog.Println("Client initialized...")

	ticker := time.NewTicker(1 * time.Second)
	startTime := time.Now()

	for {
		tickedTime := <-ticker.C

		timeDiffSec := int64(tickedTime.Sub(startTime).Seconds())
		if timeDiffSec%int64(agentConfig.PollInterval) == 0 {
			collector.UpdateMetrics()
			metrics.InfoLog.Println("Metrics have been updated")
		}
		if timeDiffSec%int64(agentConfig.PollInterval) == 0 {
			PushMetrics(client, endpoint, collector.GetMetrics(), agentConfig.Key)
			metrics.InfoLog.Println("Metrics have been pushed")
			GetMetricsValues(client, endpoint, agentConfig.Key, collector.GetMetrics())
			metrics.InfoLog.Println("Metrics update has been received")
		}
	}
}
