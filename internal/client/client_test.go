package client

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	col "github.com/nmramorov/gowatcher/internal/collector"
	"github.com/nmramorov/gowatcher/internal/config"
	"github.com/nmramorov/gowatcher/internal/log"
	pb "github.com/nmramorov/gowatcher/internal/proto"
)

func TestPushMetrics(t *testing.T) {
	collector := col.NewCollector()
	endpoint := "http://127.0.0.1:8080"

	client := &http.Client{}
	assert.NotPanics(t, func() { PushMetrics(client, endpoint, collector.GetMetrics(), "", "") })
}

func TestCreateRequests(t *testing.T) {
	collector := col.NewCollector()
	endpoint := "http://127.0.0.1:8080"
	assert.IsType(t, make([]*http.Request, 0), CreateRequests(endpoint, collector.GetMetrics()))
}

func TestPushMetricsBatch(t *testing.T) {
	collector := col.NewCollector()
	endpoint := "http://127.0.0.1:8080"

	client := &http.Client{}
	assert.NotPanics(t, func() { PushMetricsBatch(client, endpoint, "", collector.GetMetrics()) })
}

func TestGetMetricsValues(t *testing.T) {
	collector := col.NewCollector()
	endpoint := "http://127.0.0.1:8080"

	client := &http.Client{}
	assert.NotPanics(t, func() { GetMetricsValues(client, endpoint, "gauge", "", collector.GetMetrics()) })
}

func TestClientRun(t *testing.T) {
	testFoo := func() {
		client := Client{}
		client.Run()
	}
	assert.NotPanics(t, testFoo)
}

func TestRunConcurrentlyNoGRPC(t *testing.T) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, time.Duration(5)*time.Second)
	defer cancel()

	testFoo := func() {
		clientKill := make(chan struct{}, 1)
		agentConfig := &config.AgentConfig{
			Address:        "localhost:8080",
			ReportInterval: 2,
			PollInterval:   1,
			Key:            "some key",
			RateLimit:      55,
			PublicKeyPath:  "wrong path",
			GRPC:           false,
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
			// получаем переменную интерфейсного типа UsersClient,
			// через которую будем отправлять сообщения
			grpcClient = pb.NewMetricsClient(conn)
		}

		RunConcurrently(ctx, clientKill, agentConfig, client, grpcClient, endpoint)
	}
	require.NotPanics(t, testFoo)
}

func TestRunConcurrentlyGRPC(t *testing.T) {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, time.Duration(5)*time.Second)
	defer cancel()

	testFoo := func() {
		clientKill := make(chan struct{}, 1)
		agentConfig := &config.AgentConfig{
			Address:        "localhost:3200",
			ReportInterval: 2,
			PollInterval:   1,
			Key:            "some key",
			RateLimit:      55,
			PublicKeyPath:  "wrong path",
			GRPC:           true,
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
			grpcClient = pb.NewMetricsClient(conn)
		}

		RunConcurrently(ctx, clientKill, agentConfig, client, grpcClient, endpoint)
	}
	require.NotPanics(t, testFoo)
}

