package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	col "github.com/nmramorov/gowatcher/internal/collector"
	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/config"
	"github.com/nmramorov/gowatcher/internal/hashgen"
	"github.com/nmramorov/gowatcher/internal/log"
	pb "github.com/nmramorov/gowatcher/internal/proto"
	sec "github.com/nmramorov/gowatcher/internal/security"
)

var (
	COUNTER = "counter"
	GAUGE   = "gauge"
)

func CreateRequests(endpoint string, mtrcs *m.Metrics) []*http.Request {
	requests := make([]*http.Request, 0)
	for k, v := range mtrcs.GaugeMetrics {
		req, err := http.NewRequest(http.MethodPost, endpoint+fmt.Sprintf("/update/gauge/%s/%f", k, v), http.NoBody)
		if err != nil {
			log.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "text/plain")
		req.Header.Add("X-Real-IP", endpoint)
		requests = append(requests, req)
	}
	for k, v := range mtrcs.CounterMetrics {
		req, err := http.NewRequest(http.MethodPost, endpoint+fmt.Sprintf("/update/counter/%s/%d", k, v), http.NoBody)
		if err != nil {
			log.ErrorLog.Printf("Could not do POST request for counter with params: %s %d", k, v)
		}
		req.Header.Add("Content-Type", "text/plain")
		req.Header.Add("X-Real-IP", endpoint)
		requests = append(requests, req)
	}
	return requests
}

func createBatch(src *m.Metrics) []*m.JSONMetrics {
	batchCap := len(src.CounterMetrics) + len(src.GaugeMetrics)
	batch := make([]*m.JSONMetrics, 0, batchCap)
	for k, v := range src.GaugeMetrics {
		v := v
		batch = append(batch, &m.JSONMetrics{
			ID:    k,
			MType: GAUGE,
			Value: (*float64)(&v),
		})
	}
	for k, v := range src.CounterMetrics {
		v := v
		batch = append(batch, &m.JSONMetrics{
			ID:    k,
			MType: COUNTER,
			Delta: (*int64)(&v),
		})
	}
	return batch
}

func encodeBatch(batch []*m.JSONMetrics, certPath string) *bytes.Buffer {
	body := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(body)
	err := encoder.Encode(&batch)
	if err != nil {
		log.ErrorLog.Printf("Error during batch encoding: %e", err)
	}
	cert, err := sec.GetCertificate(certPath)
	if err != nil {
		log.ErrorLog.Printf("error getting cert: %e", err)
	}
	encodedMsg, err := sec.EncodeMsg(body.Bytes(), cert)
	if err != nil {
		log.ErrorLog.Printf("error encoding batch: %e", err)
		return body
	}
	buf := bytes.NewBuffer(encodedMsg)
	return buf
}

func createRequestsBatch(endpoint, path, certPath string, src *m.Metrics) *http.Request {
	batch := createBatch(src)
	body := encodeBatch(batch, certPath)
	req, err := http.NewRequest(http.MethodPost, endpoint+path, body)
	if err != nil {
		log.ErrorLog.Println("Could not do POST batch request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Real-IP", endpoint)
	return req
}

func createBody(metricType, path, key, secretkey, certPath string, value interface{}) *bytes.Buffer {
	var hash string
	if secretkey != "" {
		generator := hashgen.NewHashGenerator(secretkey)
		hash = generator.GenerateHash(metricType, key, value)
	}
	body := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(body)
	var toEncode m.JSONMetrics
	switch metricType {
	case GAUGE:
		toEncode.MType = GAUGE
		toEncode.ID = key
		gaugeVal := value.(m.Gauge)
		val := (*float64)(&gaugeVal)
		toEncode.Value = val
		toEncode.Hash = hash
	case COUNTER:
		toEncode.MType = COUNTER
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
	err := encoder.Encode(&toEncode)
	if err != nil {
		log.ErrorLog.Printf("error encoding body: %e", err)
	}
	cert, err := sec.GetCertificate(certPath)
	if err != nil {
		log.ErrorLog.Printf("error getting cert: %e", err)
		return body
	}
	encodedMsg, _ := sec.EncodeMsg(body.Bytes(), cert)
	buf := bytes.NewBuffer(encodedMsg)
	return buf
}

func createGaugeRequests(endpoint, path, key, certPath string, gaugeMetrics map[string]m.Gauge) []*http.Request {
	requests := make([]*http.Request, 0)
	for k, v := range gaugeMetrics {
		body := createBody(GAUGE, path, k, key, certPath, v)

		req, err := http.NewRequest(http.MethodPost, endpoint+path, body)
		if err != nil {
			log.ErrorLog.Printf("Could not do POST request for gauge with params: %s %f", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Real-IP", endpoint)
		requests = append(requests, req)
	}
	return requests
}

func createCounterRequests(endpoint, path, key, certPath string, counterMetrics map[string]m.Counter) []*http.Request {
	requests := make([]*http.Request, 0)
	for k, v := range counterMetrics {
		body := createBody(COUNTER, path, k, key, certPath, v)
		req, err := http.NewRequest(http.MethodPost, endpoint+path, body)
		if err != nil {
			log.ErrorLog.Printf("Could not do POST request for counter with params: %s %d", k, v)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Real-IP", endpoint)
		requests = append(requests, req)
	}
	return requests
}

func generateMetricsRequests(endpoint, path, key, certPath string, src *m.Metrics) []*http.Request {
	gaugeRequests := createGaugeRequests(endpoint, path, key, certPath, src.GaugeMetrics)
	counterRequests := createCounterRequests(endpoint, path, key, certPath, src.CounterMetrics)
	return append(gaugeRequests, counterRequests...)
}

func PushMetrics(client *http.Client, endpoint string, mtrcs *m.Metrics, key, certPath string) {
	defer func() {
		if p := recover(); p != nil {
			log.ErrorLog.Println(p)
		}
	}()
	requests := generateMetricsRequests(endpoint, "/update/", key, certPath, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			log.ErrorLog.Println(err)
		}
		err = resp.Body.Close()
		if err != nil {
			log.ErrorLog.Printf("error during response body close: %e", err)
		}
	}
}

func PushMetricsBatch(client *http.Client, endpoint, certPath string, mtrcs *m.Metrics) {
	defer func() {
		if p := recover(); p != nil {
			log.ErrorLog.Println(p)
		}
	}()
	request := createRequestsBatch(endpoint, "/updates/", certPath, mtrcs)
	resp, err := client.Do(request)
	if err != nil {
		log.ErrorLog.Println(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.ErrorLog.Printf("error during response body close: %e", err)
	}
}

func GetMetricsValues(client *http.Client, endpoint, key, certPath string, mtrcs *m.Metrics) {
	defer func() {
		if p := recover(); p != nil {
			log.ErrorLog.Println(p)
		}
	}()
	requests := generateMetricsRequests(endpoint, "/value/", key, certPath, mtrcs)
	for _, request := range requests {
		resp, err := client.Do(request)
		if err != nil {
			log.ErrorLog.Println(err)
		}
		err = resp.Body.Close()
		if err != nil {
			log.ErrorLog.Printf("error during response body close: %e", err)
		}
	}
}

type Job struct {
	RequestType string
}

func RunTickers(ctx context.Context, stateSig chan struct{}, agentConfig *config.AgentConfig, jobCh chan<- *Job) {
	log.InfoLog.Println("Tickers are running...")
	updateTicker := time.NewTicker(time.Duration(agentConfig.PollInterval) * time.Second)
	pushTicker := time.NewTicker(time.Duration(agentConfig.ReportInterval) * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.InfoLog.Println("received context cancel func, shutting down tickers")
			close(jobCh)
			return
		case <-stateSig:
			log.InfoLog.Println("received shutdown signal, shutting down tickers")
			close(jobCh)

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

func RunConcurrently(ctx context.Context, stateSignal chan struct{}, config *config.AgentConfig,
	client *http.Client, grpcClient pb.MetricsClient, endpoint string,
) {
	jobCh := make(chan *Job, config.RateLimit)
	var wg sync.WaitGroup
	collector := col.NewCollector()
	wg.Add(1)
	go func() {
		RunTickers(ctx, stateSignal, config, jobCh)
		wg.Done()
	}()
	select {
	case <-ctx.Done():
		wg.Wait()
		return
	default:
		for job := range jobCh {
			log.InfoLog.Printf("Running job %s", job.RequestType)
			RunJob(ctx, job, collector, client, grpcClient, endpoint, config)
		}
	}
	wg.Wait()
}

func PushMetricsGRPC(ctx context.Context, col *col.Collector, client pb.MetricsClient) {
	for k, v := range col.Metrics.CounterMetrics {
		resp, err := client.AddMetric(ctx, &pb.AddMetricRequest{
			Metric: &pb.Metric{
				Id:    k,
				Mtype: "counter",
				Delta: int64(v),
			},
		})
		if err != nil {
			log.ErrorLog.Printf("GRPC. Error pushing metric %s: %e", k, err)
		}
		if resp != nil && resp.Error != "" {
			log.ErrorLog.Printf("GRPC resp error: %s", resp.Error)
		}
	}
	for k, v := range col.Metrics.GaugeMetrics {
		resp, err := client.AddMetric(ctx, &pb.AddMetricRequest{
			Metric: &pb.Metric{
				Id:    k,
				Mtype: "counter",
				Value: float64(v),
			},
		})
		if err != nil {
			log.ErrorLog.Printf("GRPC. Error pushing metric %s: %e", k, err)
		}
		if resp != nil && resp.Error != "" {
			log.ErrorLog.Printf("GRPC resp error: %s", resp.Error)
		}
	}
}

func GetMetricsGRPC(ctx context.Context, metrics *m.Metrics, client pb.MetricsClient) {
	for k := range metrics.CounterMetrics {
		resp, err := client.GetMetric(
			ctx, &pb.GetMetricRequest{
				Metric: &pb.Metric{
					Id:    k,
					Mtype: "counter",
				},
			},
		)
		if err != nil {
			log.ErrorLog.Printf("GRPC. Error getting metric %s: %e", k, err)
		} else {
			if resp.Error != "" {
				log.ErrorLog.Printf("GRPC resp error: %s", resp.Error)
			} else {
				log.InfoLog.Printf("GRPC received metric %s:", resp.Metric)
			}
		}
	}
	for k := range metrics.GaugeMetrics {
		resp, err := client.GetMetric(
			ctx, &pb.GetMetricRequest{
				Metric: &pb.Metric{
					Id:    k,
					Mtype: "gauge",
				},
			},
		)
		if err != nil {
			log.ErrorLog.Printf("GRPC. Error getting metric %s: %e", k, err)
		} else {
			if resp.Error != "" {
				log.ErrorLog.Printf("GRPC resp error: %s", resp.Error)
			} else {
				log.InfoLog.Printf("GRPC received metric %s:", resp.Metric)
			}
		}
	}
}

func RunJob(ctx context.Context, job *Job, collector *col.Collector, client *http.Client,
	grpcClient pb.MetricsClient, endpoint string, agentConfig *config.AgentConfig,
) {
	jobCtx, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		log.InfoLog.Println("Context cancel was called, stopping job")
		return
	default:
		switch job.RequestType {
		case "update":
			collector.UpdateMetrics()
			log.InfoLog.Println("Metrics have been updated concurrently")
			collector.UpdateExtraMetrics()
			log.InfoLog.Println("Extra metrics have been updated concurrently")
		case "push":
			if agentConfig.GRPC {
				PushMetricsGRPC(
					jobCtx, collector, grpcClient,
				)
				log.InfoLog.Println("Metrics have been pushed via GRPC")
				GetMetricsGRPC(jobCtx, collector.GetMetrics(), grpcClient)
				log.InfoLog.Println("Metrics have been received via GRPC")
			} else {
				PushMetrics(client, endpoint, collector.GetMetrics(), agentConfig.Key, agentConfig.PublicKeyPath)
				log.InfoLog.Println("Metrics have been pushed")
				PushMetricsBatch(client, endpoint, agentConfig.PublicKeyPath, collector.GetMetrics())
				log.InfoLog.Println("Batch metrics were pushed")
				GetMetricsValues(client, endpoint, agentConfig.Key, agentConfig.PublicKeyPath, collector.GetMetrics())
				log.InfoLog.Println("Metrics update has been received")
			}
		}
	}
}

type Client struct{}

func (c *Client) Run() {
	ctx := context.Background()
	idleConnsClosed := make(chan struct{})

	sigint := make(chan os.Signal, 1)
	clientKill := make(chan struct{}, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigint
		close(clientKill)
		close(idleConnsClosed)
	}()

	agentConfig, err := config.GetAgentConfig()
	if err != nil {
		log.ErrorLog.Printf("Error with agent config: %e", err)
		return
	}
	endpoint := "http://" + agentConfig.Address

	client := &http.Client{}
	var grpcClient pb.MetricsClient
	if agentConfig.GRPC {
		conn, err := grpc.Dial(":"+strings.Split(agentConfig.Address, ":")[1],
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.ErrorLog.Fatal(err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				log.ErrorLog.Println("error closing grpc client")
			}
		}()
		// получаем переменную интерфейсного типа MetricsClient,
		// через которую будем отправлять сообщения
		grpcClient = pb.NewMetricsClient(conn)
	}

	RunConcurrently(ctx, clientKill, agentConfig, client, grpcClient, endpoint)
	<-idleConnsClosed
}
