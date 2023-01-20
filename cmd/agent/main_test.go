package main

import (
	"net/http"
	"testing"

	"internal/metrics"

	"github.com/stretchr/testify/assert"
)

func TestPushMetrics(t *testing.T) {
	var collector = metrics.NewCollector()
	endpoint := "http://127.0.0.1:8080"

	client := &http.Client{}
	assert.NotPanics(t, func() { PushMetrics(client, endpoint, collector.GetMetrics()) })
}

func TestCreateRequests(t *testing.T) {
	var collector = metrics.NewCollector()
	endpoint := "http://127.0.0.1:8080"
	assert.IsType(t, make([]*http.Request, 0), CreateRequests(endpoint, collector.GetMetrics()))
}
