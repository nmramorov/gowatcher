package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	col "github.com/nmramorov/gowatcher/internal/collector"
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
