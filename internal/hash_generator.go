package metrics

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type HashGenerator struct {
	secretkey string
}

func (gen *HashGenerator) GenerateHash(metricType, id string, value interface{}) string {
	var hashString string
	switch metricType {
	case "gauge":
		hashString = fmt.Sprintf("%s:gauge:%f", id, value) + gen.secretkey
	case "counter":
		hashString = fmt.Sprintf("%s:counter:%d", id, value) + gen.secretkey
	}
	h := sha256.New()
	h.Write([]byte(hashString))
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}

func NewHashGenerator(key string) *HashGenerator {
	return &HashGenerator{
		secretkey: key,
	}
}
