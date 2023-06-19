package file

import (
	"encoding/json"
	"os"

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/log"
)

type FileReader struct {
	file    *os.File
	decoder *json.Decoder
}

type FileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewFileReader(fileName string) (*FileReader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &FileReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (fr *FileReader) ReadJSON() (*metrics.Metrics, error) {
	metric := &metrics.Metrics{}
	if err := fr.decoder.Decode(&metric); err != nil {
		return nil, err
	}
	log.InfoLog.Println(metric)
	return metric, nil
}

func (fr *FileReader) Close() error {
	return fr.file.Close()
}

func NewFileWriter(fileName string) (*FileWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &FileWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (fw *FileWriter) WriteJSON(metric *metrics.Metrics) error {
	return fw.encoder.Encode(&metric)
}

func (fw *FileWriter) Close() error {
	return fw.file.Close()
}
