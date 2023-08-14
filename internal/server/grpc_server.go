package server

import (
	"context"
	"fmt"

	// импортируем пакет со сгенерированными protobuf-файлами
	"github.com/nmramorov/gowatcher/internal/api/handlers"
	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/log"
	pb "github.com/nmramorov/gowatcher/internal/proto"
)

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer

	h handlers.Handler
}

// AddMetric реализует интерфейс добавления метрики.
func (s *MetricsServer) AddMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse
	metricToAdd := m.JSONMetrics{
		ID:    in.Metric.Id,
		MType: in.Metric.Mtype,
		Delta: &in.Metric.Delta,
		Value: &in.Metric.Value,
		Hash:  in.Metric.Hash,
	}
	updatedData, err := s.h.Collector.UpdateMetricFromJSON(&metricToAdd)
	if err != nil {
		log.ErrorLog.Printf("Error occurred during metric update from json: %e", err)
		response.Error = fmt.Sprintf("Error occurred during metric update from json: %e", err)
	}
	if s.h.Cursor.IsValid {
		err = s.h.Cursor.Add(ctx, updatedData)
		if err != nil {
			log.ErrorLog.Printf("could not add data to db for metric %s: %e", metricToAdd.ID, err)
			response.Error = fmt.Sprintf("could not add data to db for metric %s: %e", metricToAdd.ID, err)
		}
	}
	return &response, nil
}

// GetMetric реализует интерфейс получения метрики.
func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var response pb.GetMetricResponse
	metricToAdd := m.JSONMetrics{
		ID:    in.Metric.Id,
		MType: in.Metric.Mtype,
	}
	var metric *m.JSONMetrics
	var err error
	if s.h.Cursor.IsValid {
		metric, err = s.h.Cursor.Get(ctx, &metricToAdd)
		if err != nil {
			log.ErrorLog.Printf("could not get data from db for metric %s: %e", metricToAdd.ID, err)
			response.Error = fmt.Sprintf("could not get data from db for metric %s: %e", metricToAdd.ID, err)
		}
	}
	if metric == nil {
		metric, err = s.h.Collector.GetMetricJSON(&metricToAdd)
		if err != nil {
			log.ErrorLog.Printf("Error occurred during metric getting from json: %e", err)
			response.Error = fmt.Sprintf("Error occurred during metric getting from json: %e", err)
		}
	}
	response.Metric = &pb.Metric{
		Id:    metric.ID,
		Mtype: metric.MType,
		Delta: *metric.Delta,
		Value: *metric.Value,
		Hash:  metric.Hash,
	}

	return &response, nil
}
