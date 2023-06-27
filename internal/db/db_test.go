package db

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/jackc/pgx/v5/stdlib"
	mock_db "github.com/nmramorov/gowatcher/internal/db/mocks"
	"github.com/stretchr/testify/require"

	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/errors"
)

func TestNewCursor(t *testing.T) {
	parent := context.Background()
	type args struct {
		parent  context.Context
		link    string
		adaptor string
	}
	type test struct {
		name    string
		args    args
		wantErr bool
	}
	tests := []test{
		{
			name: "Positive Cursor creation",
			args: args{
				parent:  parent,
				link:    "",
				adaptor: "pgx",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCursor(tt.args.parent, tt.args.link, tt.args.adaptor)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCursor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCursor_Close(t *testing.T) {
	parent := context.Background()

	type args struct {
		parent context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Close Positive case",
			args: args{
				parent: parent,
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCursor(tt.args.parent, "", "pgx")
			if err != nil {
				t.Errorf("Error creating Cursror %v", err)
			}
			if err := c.CloseConnection(tt.args.parent); (err != nil) != tt.wantErr {
				t.Errorf("Cursor.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCursor_Ping(t *testing.T) {
	parent := context.Background()

	cursor, _ := NewCursor(parent, "", "pgx")

	type args struct {
		parent context.Context
	}
	tests := []struct {
		name    string
		c       *Cursor
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Ping Negative",
			c:    cursor,
			args: args{
				parent: parent,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Ping(tt.args.parent); (err != nil) != tt.wantErr {
				t.Errorf("Cursor.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitDBPositive(t *testing.T) {
	parent := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_db.NewMockDriverMethods(ctrl)
	s.EXPECT().ExecContext(gomock.Any(), CreateGaugeTable).Return(nil, nil).MaxTimes(1)
	s.EXPECT().ExecContext(gomock.Any(), CreateCounterTable).Return(nil, nil).MaxTimes(1)
	c := Cursor{
		DB: s,
	}
	require.NoError(t, c.InitDB(parent))
}

func TestInitDBNegative(t *testing.T) {
	parent := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_db.NewMockDriverMethods(ctrl)
	s.EXPECT().ExecContext(gomock.Any(), CreateGaugeTable).Return(nil, errors.ErrorMetricNotFound).MaxTimes(1)
	s.EXPECT().ExecContext(gomock.Any(), CreateCounterTable).Return(nil, errors.ErrorMetricNotFound).MaxTimes(1)
	c := Cursor{
		DB: s,
	}
	require.Error(t, c.InitDB(parent))

	ss := mock_db.NewMockDriverMethods(ctrl)
	ss.EXPECT().ExecContext(gomock.Any(), CreateGaugeTable).Return(nil, nil).MaxTimes(1)
	ss.EXPECT().ExecContext(gomock.Any(), CreateCounterTable).Return(nil, errors.ErrorMetricNotFound).MaxTimes(1)
	c = Cursor{
		DB: ss,
	}
	require.Error(t, c.InitDB(parent))
}

func TestAddPositive(t *testing.T) {
	parent := context.Background()
	mockVal := 55.3
	mockDelta := int64(3)
	mockGaugeMetric := &m.JSONMetrics{
		ID:    "1",
		MType: "gauge",
		Value: &mockVal,
	}
	mockCounterMetric := &m.JSONMetrics{
		ID:    "1",
		MType: "counter",
		Delta: &mockDelta,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_db.NewMockDriverMethods(ctrl)
	s.EXPECT().
		ExecContext(gomock.Any(), InsertIntoGauge, mockGaugeMetric.ID, mockGaugeMetric.MType, mockGaugeMetric.Value).
		MaxTimes(1)
	c := Cursor{
		DB: s,
	}
	require.NoError(t, c.Add(parent, mockGaugeMetric))

	s = mock_db.NewMockDriverMethods(ctrl)
	s.EXPECT().
		ExecContext(gomock.Any(), InsertIntoCounter, mockCounterMetric.ID, mockCounterMetric.MType, mockCounterMetric.Delta).
		MaxTimes(1)
	c = Cursor{
		DB: s,
	}
	require.NoError(t, c.Add(parent, mockCounterMetric))
}

func TestAddNegative(t *testing.T) {
	parent := context.Background()
	mockVal := 55.3
	mockDelta := int64(3)
	mockGaugeMetric := &m.JSONMetrics{
		ID:    "1",
		MType: "gauge",
		Value: &mockVal,
	}
	mockCounterMetric := &m.JSONMetrics{
		ID:    "1",
		MType: "counter",
		Delta: &mockDelta,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_db.NewMockDriverMethods(ctrl)
	s.EXPECT().
		ExecContext(gomock.Any(), InsertIntoGauge, mockGaugeMetric.ID, mockGaugeMetric.MType, mockGaugeMetric.Value).
		Return(nil, errors.ErrorWithIntervalConvertion).
		MaxTimes(1)
	c := Cursor{
		DB: s,
	}
	require.Error(t, c.Add(parent, mockGaugeMetric))

	s = mock_db.NewMockDriverMethods(ctrl)
	s.EXPECT().
		ExecContext(gomock.Any(), InsertIntoCounter, mockCounterMetric.ID, mockCounterMetric.MType, mockCounterMetric.Delta).
		Return(nil, errors.ErrorWithIntervalConvertion).
		MaxTimes(1)
	c = Cursor{
		DB: s,
	}
	require.Error(t, c.Add(parent, mockCounterMetric))
}

func TestAddBatchPositive(t * testing.T) {
	parent := context.Background()
	mockVal := 55.3
	mockGaugeMetrics := make([]*m.JSONMetrics, 0, 1)
	mockGaugeMetrics = append(mockGaugeMetrics, &m.JSONMetrics{
		ID:    "1",
		MType: "gauge",
		Value: &mockVal,
	})
	c := Cursor{
		buffer: make([]*m.JSONMetrics, 0, 3),
	}
	require.NoError(t, c.AddBatch(parent, mockGaugeMetrics))
}
