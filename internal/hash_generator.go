package metrics

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Тип, суть которого - сгенерировать хеш.
type HashGenerator struct {
	secretkey string
}

// Метод, генерирующий хеши для метрик разных типов.
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

// Конструктор для типа HashGenerator.
func NewHashGenerator(key string) *HashGenerator {
	return &HashGenerator{
		secretkey: key,
	}
}
