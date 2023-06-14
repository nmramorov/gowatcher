package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	col "github.com/nmramorov/gowatcher/internal/collector"
	"github.com/nmramorov/gowatcher/internal/log"
	"github.com/nmramorov/gowatcher/internal/hashgen"
	"github.com/nmramorov/gowatcher/internal/config"
)

func CreateRequests(endpoint string, mtrcs *m.Metrics) []*http.Request {
	var requests []*http.Request
	for k, v := range mtrcs.GaugeMetrics {
		req, err := http.NewRequest(http.MethodPost, endpoint+fmt.Sprintf("/update/gauge/%s/%f", k, v), nil)
		if err != nil {
			log.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "text/plain")
		requests = append(requests, req)
	}
	for k, v := range mtrcs.CounterMetrics {
		req, err := http.NewRequest(http.MethodPost, endpoint+fmt.Sprintf("/update/counter/%s/%d", k, v), nil)
		if err != nil {
			log.ErrorLog.Printf("Could not do POST request for counter with params: %s %d", k, v)
		}
		req.Header.Add("Content-Type", "text/plain")
		requests = append(requests, req)
	}
	return requests
}

func createBatch(src *m.Metrics) []*m.JSONMetrics {
	var batch []*m.JSONMetrics
	for k, v := range src.GaugeMetrics {
		batch = append(batch, &m.JSONMetrics{
			ID:    k,
			MType: "gauge",
			Value: (*float64)(&v),
		})
	}
	for k, v := range src.CounterMetrics {
		batch = append(batch, &m.JSONMetrics{
			ID:    k,
			MType: "counter",
			Delta: (*int64)(&v),
		})
	}
	return batch
}

func encodeBatch(batch []*m.JSONMetrics) *bytes.Buffer {
	body := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(body)
	encoder.Encode(&batch)
	return body
}

func createRequestsBatch(endpoint, path string, src *m.Metrics) *http.Request {
	batch := createBatch(src)
	body := encodeBatch(batch)
	req, err := http.NewRequest(http.MethodPost, endpoint+path, body)
	if err != nil {
		log.ErrorLog.Println("Could not do POST batch request")
	}
	req.Header.Add("Content-Type", "application/json")
	return req
}

func createBody(metricType, path, key, secretkey string, value interface{}) *bytes.Buffer {
	var hash string
	if secretkey != "" {
		generator := hashgen.NewHashGenerator(secretkey)
		hash = generator.GenerateHash(metricType, key, value)
	}
	body := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(body)
	var toEncode m.JSONMetrics
	switch metricType {
	case "gauge":
		toEncode.MType = "gauge"
		toEncode.ID = key
		gaugeVal := value.(m.Gauge)
		val := (*float64)(&gaugeVal)
		toEncode.Value = val
		toEncode.Hash = hash
	case "counter":
		toEncode.MType = "counter"
		toEncode.ID = key
		counterVal := value.(m.Counter)
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

func createGaugeRequests(endpoint, path, key string, gaugeMetrics map[string]m.Gauge) []*http.Request {
	var requests []*http.Request
	for k, v := range gaugeMetrics {
		body := createBody("gauge", path, k, key, v)
		req, err := http.NewRequest(http.MethodPost, endpoint+path, body)
		if err != nil {
			log.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	return requests
}

func createCounterRequests(endpoint, path, key string, counterMetrics map[string]m.Counter) []*http.Request {
	var requests []*http.Request
	for k, v := range counterMetrics {
		body := createBody("counter", path, k, key, v)
		req, err := http.NewRequest(http.MethodPost, endpoint+path, body)
		if err != nil {
			log.ErrorLog.Printf("Could not do POST request for counter with params: %s %d", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		requests = append(requests, req)
	}
	return requests
}

func generateMetricsRequests(endpoint, path, key string, src *m.Metrics) []*http.Request {
	gaugeRequests := createGaugeRequests(endpoint, path, key, src.GaugeMetrics)
	counterRequests := createCounterRequests(endpoint, path, key, src.CounterMetrics)
	return append(gaugeRequests, counterRequests...)
}

func PushMetrics(client *http.Client, endpoint string, mtrcs *m.Metrics, key string) {
	defer func() {
		if p := recover(); p != nil {
			log.ErrorLog.Println(p)
		}
	}()
	requests := generateMetricsRequests(endpoint, "/update/", key, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			log.ErrorLog.Println(err)
		}
		defer resp.Body.Close()
	}
}

func PushMetricsBatch(client *http.Client, endpoint string, mtrcs *m.Metrics) {
	defer func() {
		if p := recover(); p != nil {
			log.ErrorLog.Println(p)
		}
	}()
	request := createRequestsBatch(endpoint, "/updates/", mtrcs)
	resp, err := client.Do(request)
	if err != nil {
		log.ErrorLog.Println(err)
	}
	defer resp.Body.Close()
}

func GetMetricsValues(client *http.Client, endpoint, key string, mtrcs *m.Metrics) {
	defer func() {
		if p := recover(); p != nil {
			log.ErrorLog.Println(p)
		}
	}()
	requests := generateMetricsRequests(endpoint, "/value/", key, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			log.ErrorLog.Println(err)
		}
		defer resp.Body.Close()
	}
}

type Job struct {
	RequestType string
}

func RunTickers(agentConfig *config.AgentConfig, jobCh chan<- *Job) {
	log.InfoLog.Println("Tickers are running...")
	updateTicker := time.NewTicker(time.Duration(agentConfig.PollInterval) * time.Second)
	pushTicker := time.NewTicker(time.Duration(agentConfig.ReportInterval) * time.Second)

	for {
		select {
		case <-updateTicker.C:
			jobCh <- &Job{
				RequestType: "update",
			}
		case <-pushTicker.C:
			jobCh <- &Job{
				RequestType: "push",
			}
		}
	}
}

func RunConcurrently(config *config.AgentConfig, client *http.Client, endpoint string) {
	jobCh := make(chan *Job, config.RateLimit)

	var collector = col.NewCollector()
	go func() {
		RunTickers(config, jobCh)
	}()

	for job := range jobCh {
		log.InfoLog.Printf("Running job %s", job.RequestType)
		RunJob(job, collector, client, endpoint, config)
	}

}

func RunJob(job *Job, collector *col.Collector, client *http.Client, endpoint string, agentConfig *config.AgentConfig) {
	switch job.RequestType {
	case "update":
		go func() {
			collector.UpdateMetrics()
			log.InfoLog.Println("Metrics have been updated concurrently")
		}()
		go func() {
			collector.UpdateExtraMetrics()
			log.InfoLog.Println("Extra metrics have been updated concurrently")
		}()
	case "push":
		go func() {
			PushMetrics(client, endpoint, collector.GetMetrics(), agentConfig.Key)
			log.InfoLog.Println("Metrics have been pushed")
			PushMetricsBatch(client, endpoint, collector.GetMetrics())
			log.InfoLog.Println("Batch metrics were pushed")
			GetMetricsValues(client, endpoint, agentConfig.Key, collector.GetMetrics())
			log.InfoLog.Println("Metrics update has been received")
		}()
	}
}

type Client struct {}

func (c *Client) Run() {
	agentConfig := config.GetAgentConfig()
	endpoint := "http://" + agentConfig.Address

	client := &http.Client{}

	RunConcurrently(agentConfig, client, endpoint)
}
