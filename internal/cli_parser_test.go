package metrics

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerCLI(t *testing.T) {
	os.Args = []string{"main.go", "-a", "localhost:38731", "-r=true", "-i=5m", "-f=/tmp/wmSoUM", "-k=aaab", "-d=ddd"}
	var address = flag.String("a", "localhost:8080", "server address")
	var restore = flag.String("r", "default", "restore metrics from file")
	var storeInterval = flag.String("i", "300s", "period between file save")
	var storeFile = flag.String("f", "/tmp/devops-metrics-db.json", "name of file where metrics stored")
	var key = flag.String("k", "", "key to calculate hash")
	var database = flag.String("d", "", "database link")
	flag.Parse()

	args := &ServerCLIOptions{
		Address:       *address,
		Restore:       *restore,
		StoreInterval: *storeInterval,
		StoreFile:     *storeFile,
		Key:           *key,
		Database:      *database,
	}
	assert.Equal(t, "localhost:38731", args.Address)
	assert.Equal(t, "true", args.Restore)
	assert.Equal(t, "5m", args.StoreInterval)
	assert.Equal(t, "/tmp/wmSoUM", args.StoreFile)
	assert.Equal(t, "aaab", args.Key)
	assert.Equal(t, "ddd", args.Database)
}
