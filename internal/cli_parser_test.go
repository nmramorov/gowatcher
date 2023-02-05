package metrics

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerCLI(t *testing.T) {
	os.Args = []string{"main.go", "-a", "localhost:38731", "-r=true", "-i=5m", "-f=/tmp/wmSoUM"}
	var address = flag.String("a", "localhost:8080", "server address")
	var restore = flag.String("r", "default", "restore metrics from file")
	var storeInterval = flag.String("i", "300s", "period between file save")
	var storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	flag.Parse()

	args := &ServerCLIOptions{
		Address:       *address,
		Restore:       *restore,
		StoreInterval: *storeInterval,
		StoreFile:     *storeFile,
	}
	assert.Equal(t, "localhost:38731", args.Address)
	assert.Equal(t, "true", args.Restore)
	assert.Equal(t, "5m", args.StoreInterval)
	assert.Equal(t, "/tmp/wmSoUM", args.StoreFile)
}

func TestAgentCLI(t *testing.T) {
	os.Args = []string{"main.go", "-a", "localhost:4444", "-r=11s", "-p=5m", "-k=sadfsdfsdf"}
	var address = flag.String("a", "localhost:8080", "server address")
	var reportInterval = flag.String("r", "10s", "report interval time")
	var pollInterval = flag.String("p", "2s", "poll interval time")
	var key = flag.String("k", "", "key to calculate hash")
	flag.Parse()

	args := &AgentCLIOptions{
		Address:        *address,
		ReportInterval: *reportInterval,
		PollInterval:   *pollInterval,
		Key:            *key,
	}

	assert.Equal(t, "localhost:4444", args.Address)
	assert.Equal(t, "11s", args.ReportInterval)
	assert.Equal(t, 300, func() int {
		poll := args.GetNumericInterval("PollInterval")
		return int(poll)
	}())
	assert.Equal(t, "sadfsdfsdf", args.Key)
}

func TestServerCLIDefaults(t *testing.T) {
	var address = flag.String("a", "localhost:8080", "server address")
	var restore = flag.String("r", "default", "restore metrics from file")
	var storeInterval = flag.String("i", "300s", "period between file save")
	var storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	flag.Parse()

	args := &ServerCLIOptions{
		Address:       *address,
		Restore:       *restore,
		StoreInterval: *storeInterval,
		StoreFile:     *storeFile,
	}
	assert.Equal(t, "localhost:8080", args.Address)
	assert.Equal(t, "default", args.Restore)
	assert.Equal(t, "/tmp/devops-metrics-db.json", args.StoreFile)
	assert.Equal(t, "300s", args.StoreInterval)
}

func TestAgentCLIDefaults(t *testing.T) {
	var address = flag.String("a", "localhost:8080", "server address")
	var reportInterval = flag.String("r", "10s", "report interval time")
	var pollInterval = flag.String("p", "2s", "poll interval time")
	var key = flag.String("k", "", "key to calculate hash")
	flag.Parse()
	args := &AgentCLIOptions{
		Address:        *address,
		ReportInterval: *reportInterval,
		PollInterval:   *pollInterval,
		Key: *key,
	}
	assert.Equal(t, "localhost:8080", args.Address)
	assert.Equal(t, "10s", args.ReportInterval)
	assert.Equal(t, "2s", args.PollInterval)
	assert.Equal(t, "", args.Key)
}
