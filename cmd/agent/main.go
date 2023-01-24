package main

import (
	"bytes"
	"encoding/json"
	"flag"
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
		body := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(body)
		val := (*float64)(&v)
		toEncode := metrics.JSONMetrics{
			ID:    k,
			MType: "gauge",
			Value: val,
		}
		encoder.Encode(&toEncode)
		req, err := http.NewRequest(http.MethodPost, endpoint+"/update/", body)
		if err != nil {
			metrics.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	for k, v := range mtrcs.CounterMetrics {
		body := bytes.NewBuffer([]byte{})
		encoder := json.NewEncoder(body)
		val := (*int64)(&v)
		toEncode := metrics.JSONMetrics{
			ID:    k,
			MType: "counter",
			Delta: val,
		}
		encoder.Encode(&toEncode)
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
		encoder.Encode(&metrics.JSONMetrics{
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
		encoder.Encode(&metrics.JSONMetrics{
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
	defer func() {
		if p := recover(); p != nil {
			metrics.ErrorLog.Println(p)
		}
	}()
	requests := CreateUpdateRequestsJSON(endpoint, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			metrics.ErrorLog.Println(err)
			panic(1)
		}
		defer resp.Body.Close()
	}
}

func GetMetricsValues(client *http.Client, endpoint string, mtrcs *metrics.Metrics) {
	defer func() {
		if p := recover(); p != nil {
			metrics.ErrorLog.Println(p)
		}
	}()
	requests := CreateValueRequestsJSON(endpoint, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			metrics.ErrorLog.Println(err)
			panic(1)
		}
		defer resp.Body.Close()
	}
}

type AgentConfig struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func GetAgentConfig(config *metrics.EnvConfig, args *metrics.AgentCLIOptions) *AgentConfig {
	agentConfig := AgentConfig{}
	if config.Address == "" {
		agentConfig.Address = args.Address
	} else {
		agentConfig.Address = config.Address
	}
	if config.PollInterval == "" {
		agentConfig.PollInterval = func() int {
			poll, err := args.GetNumericInterval("PollInterval")
			if err != nil {
				panic(err)
			}
			return int(poll)
		}()
	} else {
		agentConfig.PollInterval = func() int {
			poll, err := config.GetNumericInterval("PollInterval")
			if err != nil {
				panic(err)
			}
			return int(poll)
		}()
	}
	if config.ReportInterval == "" {
		agentConfig.ReportInterval = func() int {
			rep, err := args.GetNumericInterval("ReportInterval")
			if err != nil {
				panic(err)
			}
			return int(rep)
		}()
	} else {
		agentConfig.ReportInterval = func() int {
			rep, err := config.GetNumericInterval("ReportInterval")
			if err != nil {
				panic(err)
			}
			return int(rep)
		}()
	}
	return &agentConfig
}

func main() {
	config, err := metrics.NewConfig()
	if err != nil {
		panic(err)
	}
	var address = flag.String("a", "localhost:8080", "server address")
	var reportInterval = flag.String("r", "10s", "report interval time")
	var pollInterval = flag.String("p", "2s", "poll interval time")
	flag.Parse()
	args := &metrics.AgentCLIOptions{
		Address:        *address,
		ReportInterval: *reportInterval,
		PollInterval:   *pollInterval,
	}

	agentConfig := GetAgentConfig(config, args)

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
			fmt.Println(collector.GetMetrics().CounterMetrics)
			metrics.InfoLog.Println("Metrics have been updated")
		}
		if timeDiffSec%int64(agentConfig.PollInterval) == 0 {
			PushMetrics(client, endpoint, collector.GetMetrics())
			metrics.InfoLog.Println("Metrics have been pushed")
			GetMetricsValues(client, endpoint, collector.GetMetrics())
			metrics.InfoLog.Println("Metrics update has been received")
		}
	}
}
