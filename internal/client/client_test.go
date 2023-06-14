package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	col "github.com/nmramorov/gowatcher/internal/collector"
)

func TestPushMetrics(t *testing.T) {
	var collector = col.NewCollector()
	endpoint := "http://127.0.0.1:8080"

	client := &http.Client{}
	assert.NotPanics(t, func() { PushMetrics(client, endpoint, collector.GetMetrics(), "") })
}

func TestCreateRequests(t *testing.T) {
	var collector = col.NewCollector()
	endpoint := "http://127.0.0.1:8080"
	assert.IsType(t, make([]*http.Request, 0), CreateRequests(endpoint, collector.GetMetrics()))
}
