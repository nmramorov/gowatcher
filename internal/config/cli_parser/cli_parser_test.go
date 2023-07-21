package cliparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMultiplier(t *testing.T) {
	assert.Equal(t, int64(60), getMultiplier("1m"))
	assert.Equal(t, int64(1), getMultiplier("6s"))
}

func TestPositiveNewServerCLIOptions(t *testing.T) {
	os.Args = []string{"main.go", "-a", "localhost:38731", "-r=true", "-i=5m", "-f=/tmp/wmSoUM", "-k=aaab", "-d=ddd"}
	config, err := NewServerCliOptions()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:38731", config.Address)
	assert.Equal(t, "true", config.Restore)
	assert.Equal(t, "5m", config.StoreInterval)
	assert.Equal(t, "/tmp/wmSoUM", config.StoreFile)
	assert.Equal(t, "aaab", config.Key)
	assert.Equal(t, "ddd", config.Database)

	assert.Equal(t, int64(300), config.GetNumericInterval("StoreInterval"))
	assert.Equal(t, int64(0), config.GetNumericInterval("MyInterval"))
}


func TestNegativeNewServerCLIOptions(t *testing.T) {
	os.Args = []string{"main.go", "-b", "localhost:38731", "-r=true", "-i=5m", "-f=/tmp/wmSoUM", "-k=aaab", "-d=ddd"}
	_, err := NewServerCliOptions()
	assert.Error(t, err)
}
