package metrics

import (
	"crypto/hmac"
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
		hashString = fmt.Sprintf("%s:gauge:%f", id, value)
	case "counter":
		hashString = fmt.Sprintf("%s:counter:%d", id, value)
	}
	h := hmac.New(sha256.New, []byte(gen.secretkey))
	h.Write([]byte(hashString))
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}

func NewHashGenerator(key string) *HashGenerator {
	return &HashGenerator{
		secretkey: key,
	}
}
