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

func CreateUpdateRequestsJSON(endpoint string, mtrcs *metrics.Metrics) []*http.Request {
	var requests []*http.Request
	for k, v := range mtrcs.GaugeMetrics {
		metrics.InfoLog.Println(k, v)
		body := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(body)
		val := (*float64)(&v)
		toEncode := metrics.JSONMetrics{
			ID:    k,
			MType: "gauge",
			Value: val,
		}
		encoder.Encode(toEncode)
		req, err := http.NewRequest(http.MethodPost, endpoint+"/update/", body)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	for k, v := range mtrcs.CounterMetrics {
		metrics.InfoLog.Println(k, v)
		body := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(body)
		val := (*int64)(&v)
		toEncode := metrics.JSONMetrics{
			ID:    k,
			MType: "counter",
			Delta: val,
		}
		encoder.Encode(toEncode)
		req, err := http.NewRequest(http.MethodPost, endpoint+"/update/", body)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for counter with params: %s %d", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	return requests
}

func CreateValueRequestsJSON(endpoint string, mtrcs *metrics.Metrics) []*http.Request {
	var requests []*http.Request
	for k, v := range mtrcs.GaugeMetrics {
		body := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(body)
		encoder.Encode(metrics.JSONMetrics{
			ID:    k,
			MType: "gauge",
		})
		req, err := http.NewRequest(http.MethodPost, endpoint+"/value/", body)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	for k, v := range mtrcs.CounterMetrics {
		body := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(body)
		encoder.Encode(metrics.JSONMetrics{
			ID:    k,
			MType: "counter",
		})
		req, err := http.NewRequest(http.MethodPost, endpoint+"/value/", body)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for counter with params: %s %d", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	return requests
}

func PushMetrics(client *http.Client, endpoint string, mtrcs *metrics.Metrics) {
	requests := CreateUpdateRequestsJSON(endpoint, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			metrics.ErrorLog.Println(err)
			// panic(1)
		}
		defer resp.Body.Close()
	}
}

func GetMetricsValues(client *http.Client, endpoint string, mtrcs *metrics.Metrics) {
	requests := CreateValueRequestsJSON(endpoint, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			metrics.ErrorLog.Println(err)
			// panic(1)
		}
		defer resp.Body.Close()
	}
}

func main() {
	var pollInterval = 2
	var reportInterval = 10

	var collector = metrics.NewCollector()

	endpoint := "http://127.0.0.1:8080"

	client := &http.Client{}
	metrics.InfoLog.Println("Client initialized...")

	ticker := time.NewTicker(1 * time.Second)
	startTime := time.Now()

	for {
		tickedTime := <-ticker.C

		timeDiffSec := int(tickedTime.Sub(startTime).Seconds())
		if timeDiffSec%pollInterval == 0 {
			collector.UpdateMetrics()
			fmt.Println(collector.GetMetrics().CounterMetrics)
			metrics.InfoLog.Println("Metrics have been updated")
		}
		if timeDiffSec%reportInterval == 0 {
			PushMetrics(client, endpoint, collector.GetMetrics())
			metrics.InfoLog.Println("Metrics have been pushed")
			GetMetricsValues(client, endpoint, collector.GetMetrics())
			metrics.InfoLog.Println("Metrics update has been received")
		}
	}

}
