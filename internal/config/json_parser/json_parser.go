package jsonparser

import (
	"encoding/json"
	"os"

	"github.com/nmramorov/gowatcher/internal/log"
)

type ServerJSONConfig struct {
	Address        string `json:"address"`
	Restore        bool   `json:"restore"`
	StoreInterval  string `json:"store_interval"`
	StoreFile      string `json:"store_file"`
	Key            string `json:"key,omitempty"`
	Database       string `json:"database_dsn"`
	PrivateKeyPath string `json:"crypto_key"`
}

type AgentJSONConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
	Key            string `json:"key,omitempty"`
	RateLimit      int    `json:"rate_limit,omitempty"`
	PublicKeyPath  string `json:"crypto_key"`
}

func ReadJSONConfig[T ServerJSONConfig | AgentJSONConfig](path string) (*T, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.ErrorLog.Printf("Error reading file %s: %e", path, err)
		return nil, err
	}
	var newConfig T
	err = json.Unmarshal(content, &newConfig)
	if err != nil {
		log.ErrorLog.Printf("Error unmarshalling json config: %e", err)
		return nil, err
	}
	return &newConfig, nil
}
