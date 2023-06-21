package file

import (
	"encoding/json"
	"os"

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/log"
)

type Reader struct {
	file    *os.File
	decoder *json.Decoder
}

type Writer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0o777)
	if err != nil {
		return nil, err
	}
	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (fr *Reader) ReadJSON() (*metrics.Metrics, error) {
	metric := &metrics.Metrics{}
	if err := fr.decoder.Decode(&metric); err != nil {
		return nil, err
	}
	log.InfoLog.Println(metric)
	return metric, nil
}

func (fr *Reader) Close() error {
	return fr.file.Close()
}

func NewWriter(fileName string) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0o777)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (fw *Writer) WriteJSON(metric *metrics.Metrics) error {
	return fw.encoder.Encode(&metric)
}

func (fw *Writer) Close() error {
	return fw.file.Close()
}
