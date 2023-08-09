package cliparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMultiplier(t *testing.T) {
	assert.Equal(t, int64(60), GetMultiplier("1m"))
	assert.Equal(t, int64(1), GetMultiplier("6s"))
}

func TestPositiveNewServerCLIOptions(t *testing.T) {
	os.Args = []string{
		"main.go", "-a", "localhost:38731", "-r=true", "-i=5m",
		"-f=/tmp/wmSoUM", "-k=aaab", "-d=ddd", "-crypto-key=sfsdfsdfsd", "-c=/path/to/json",
		"-t=255.255.255.0", "-grpc=true",
	}
	config, err := NewServerCliOptions()
	assert.NoError(t, err)
	assert.Equal(t, "localhost:38731", config.Address)
	assert.Equal(t, "true", config.Restore)
	assert.Equal(t, "5m", config.StoreInterval)
	assert.Equal(t, "/tmp/wmSoUM", config.StoreFile)
	assert.Equal(t, "aaab", config.Key)
	assert.Equal(t, "ddd", config.Database)
	assert.Equal(t, "sfsdfsdfsd", config.CryptoKey)
	assert.Equal(t, "/path/to/json", config.Config)
	assert.Equal(t, "255.255.255.0", config.TrustedSubnet)
	assert.Equal(t, true, config.GRPC)

	assert.Equal(t, int64(300), config.GetNumericInterval("StoreInterval"))
	assert.Equal(t, int64(0), config.GetNumericInterval("MyInterval"))
}

func TestNegativeNewServerCLIOptions(t *testing.T) {
	os.Args = []string{
		"main.go", "-b", "localhost:38731", "-r=true",
		"-i=5m", "-f=/tmp/wmSoUM", "-k=aaab", "-d=ddd", "-t=fsdf",
	}
	_, err := NewServerCliOptions()
	assert.Error(t, err)
}

func TestPositiveNewAgentCLIOptions(t *testing.T) {
	os.Args = []string{
		"main.go", "-a", "localhost:38731", "-r=5s", "-p=5m",
		"-k=salt", "-l=13", "-crypto-key=sfsdfsdfsd", "-c=/path/to/config", "-grpc=true",
	}
	config, err := NewAgentCliOptions()
	assert.NoError(t, err)

	assert.Equal(t, int64(5), config.GetNumericInterval("ReportInterval"))
	assert.Equal(t, int64(300), config.GetNumericInterval("PollInterval"))
	assert.Equal(t, int64(0), config.GetNumericInterval("SomeInterval"))
}

func TestNegativeNewAgentCLIOptions(t *testing.T) {
	os.Args = []string{
		"main.go", "-b", "localhost:38731", "-r=true", "-i=5m",
		"-f=/tmp/wmSoUM", "-k=aaab", "-d=ddd",
	}
	_, err := NewAgentCliOptions()
	assert.Error(t, err)
}
