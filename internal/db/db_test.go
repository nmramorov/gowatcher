package db

import (
	"context"
	"database/sql"
	"fmt"

	// "database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/jackc/pgx/v5/stdlib"

	// "github.com/nmramorov/gowatcher/internal/collector/metrics"
	m "github.com/nmramorov/gowatcher/internal/collector/metrics"
	mock_db "github.com/nmramorov/gowatcher/internal/db/mocks"
	"github.com/nmramorov/gowatcher/internal/errors"
	"github.com/nmramorov/gowatcher/internal/log"
	"github.com/stretchr/testify/require"
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

func TestAddBatchV2(t *testing.T) {
	parent := context.Background()
	mockVal := 55.3
	// mockDelta := int64(3)
	mockGaugeMetric := &m.JSONMetrics{
		ID:    "1",
		MType: "gauge",
		Value: &mockVal,
	}
	// mockCounterMetric := &m.JSONMetrics{
	// 	ID:    "1",
	// 	MType: "counter",
	// 	Delta: &mockDelta,
	// }
	mockMetrics := make([]*m.JSONMetrics, 0)
	mockMetrics = append(mockMetrics, mockGaugeMetric)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_db.NewMockDriverMethods(ctrl)
	s.EXPECT().
		ExecContext(gomock.Any(), InsertIntoGauge, mockGaugeMetric.ID, mockGaugeMetric.MType, mockGaugeMetric.Value).
		MaxTimes(1)
	c := Cursor{
		DB:     s,
		buffer: make([]*m.JSONMetrics, 0),
	}
	log.InfoLog.Println(mockMetrics)
	require.NoError(t, c.AddBatchV2(parent, mockMetrics))

	// s = mock_db.NewMockDriverMethods(ctrl)
	// s.EXPECT().
	// 	ExecContext(gomock.Any(), InsertIntoCounter, mockCounterMetric.ID, mockCounterMetric.MType, mockCounterMetric.Delta).
	// 	MaxTimes(1)
	// c = Cursor{
	// 	DB: s,
	// }
	// require.NoError(t, c.Add(parent, mockCounterMetric))
}

func TestGetNegative(t *testing.T) {
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
		QueryRowContext(gomock.Any(), SelectFromGauge, mockGaugeMetric.ID).
		MaxTimes(1)
	c := Cursor{
		DB: s,
	}
	_, err := c.Get(parent, mockGaugeMetric)
	require.Error(t, err)

	s = mock_db.NewMockDriverMethods(ctrl)
	s.EXPECT().
		QueryRowContext(gomock.Any(), SelectFromCounter, mockCounterMetric.ID).
		MaxTimes(1)
	c = Cursor{
		DB: s,
	}
	_, err = c.Get(parent, mockCounterMetric)
	require.Error(t, err)

}

// func TestAddBatchPositiveNoBufCap(t *testing.T) {
// 	parent := context.Background()
// 	mockVal := 55.3
// 	mockGaugeMetrics := make([]*m.JSONMetrics, 0, 1)
// 	mockGaugeMetrics = append(mockGaugeMetrics, &m.JSONMetrics{
// 		ID:    "1",
// 		MType: "gauge",
// 		Value: &mockVal,
// 	})
// 	c := Cursor{
// 		buffer: make([]*m.JSONMetrics, 0, 3),
// 	}
// 	require.NoError(t, c.AddBatch(parent, mockGaugeMetrics))
// }

// type addMock struct {}

// func (a addMock) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {

// }

// func Test_add(t *testing.T) {
// 	parent := context.Background()

//		type args struct {
//			parent          context.Context
//			incomingMetrics *metrics.JSONMetrics
//			db              interface {
//				ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
//			}
//		}
//		tests := []struct {
//			name    string
//			parent          context.Context
//			incomingMetrics *metrics.JSONMetrics
//			db              interface {
//				ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
//			}
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//			struct{name string; parent; wantErr bool}{
//				name: "positive add",
//				parent: parent,
//				db: addMock{},
//			},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				if err := add(tt.args.parent, tt.args.incomingMetrics, tt.args.db); (err != nil) != tt.wantErr {
//					t.Errorf("add() error = %v, wantErr %v", err, tt.wantErr)
//				}
//			})
//		}
//	}
type execContextFunc func(ctx context.Context, args ...any) (sql.Result, error)

func (f execContextFunc) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	return f(ctx, args...)
}

var _ interface {
	ExecContext(ctx context.Context, args ...any) (sql.Result, error)
} = execContextFunc(nil)

type result int64

func (r result) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("unsupported")
}

func (r result) RowsAffected() (int64, error) {
	return int64(r), nil
}

func TestAzaza(t *testing.T) {
	for _, tt := range []struct {
		name string
		db   interface {
			ExecContext(ctx context.Context, args ...any) (sql.Result, error)
		}
		expectedResultRows result
		expectedError      error
	}{
		{
			name: "5 rows inserted",
			db: execContextFunc(func(ctx context.Context, args ...any) (sql.Result, error) {
				return result(5), nil
			}),
			expectedResultRows: result(5),
			expectedError:      nil,
		},
	} {
		t.Run("", func(t *testing.T) {
			r, err := tt.db.ExecContext(context.Background(), "test query")
			require.Equal(t, tt.expectedResultRows, r)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

// func TestAddBatchNegativeBufCap(t *testing.T) {
// 	parent := context.Background()
// 	mockVal := 55.3
// 	mockDelta := int64(3)
// 	mockGaugeMetric := &m.JSONMetrics{
// 		ID:    "1",
// 		MType: "gauge",
// 		Value: &mockVal,
// 	}
// 	mockCounterMetric := &m.JSONMetrics{
// 		ID:    "1",
// 		MType: "counter",
// 		Delta: &mockDelta,
// 	}
// 	mockBatch := make([]*m.JSONMetrics, 0, 2)
// 	mockBatch = append(mockBatch, mockGaugeMetric, mockCounterMetric)

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	s := mock_db.NewMockDriverMethods(ctrl)
// 	s.EXPECT().
// 		BeginTx(gomock.Any(), gomock.Any()).
// 		Return(nil, errors.ErrorWithIntervalConvertion).
// 		MaxTimes(1)
// 	c := Cursor{
// 		DB: s,
// 	}
// 	require.Error(t, c.AddBatch(parent, mockBatch))
// }

// func TestAddBatchPositiveBufCap(t *testing.T) {
// 	parent := context.Background()
// 	mockVal := 55.3
// 	mockDelta := int64(3)
// 	mockGaugeMetric := &m.JSONMetrics{
// 		ID:    "1",
// 		MType: "gauge",
// 		Value: &mockVal,
// 	}
// 	mockCounterMetric := &m.JSONMetrics{
// 		ID:    "1",
// 		MType: "counter",
// 		Delta: &mockDelta,
// 	}
// 	mockBatch := make([]*m.JSONMetrics, 0, 2)
// 	mockBatch = append(mockBatch, mockGaugeMetric, mockCounterMetric)

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	s := mock_db.NewMockDriverMethods(ctrl)
// 	s.EXPECT().
// 		BeginTx(gomock.Any(), gomock.Any()).
// 		Return(nil, errors.ErrorWithIntervalConvertion).
// 		MaxTimes(1)
// 	c := Cursor{
// 		DB: s,
// 	}
// 	require.Error(t, c.AddBatch(parent, mockBatch))
// }

// type mockTx struct {}

// func (m *mockTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
// 	return nil, errors.ErrorWithIntervalConvertion
// }

// func TestAddBatchPositiveWhenFlush(t *testing.T) {
// 	parent := context.Background()
// 	mockVal := 55.3
// 	mockDelta := int64(3)
// 	mockGaugeMetric := &m.JSONMetrics{
// 		ID:    "1",
// 		MType: "gauge",
// 		Value: &mockVal,
// 	}
// 	mockCounterMetric := &m.JSONMetrics{
// 		ID:    "1",
// 		MType: "counter",
// 		Delta: &mockDelta,
// 	}
// 	mockBatch := make([]*m.JSONMetrics, 0, 2)
// 	mockBatch = append(mockBatch, mockGaugeMetric, mockCounterMetric)

// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	s := mock_db.NewMockDriverMethods(ctrl)
// 	s.EXPECT().
// 		BeginTx(gomock.Any(), gomock.Any()).
// 		Return(&sql.Tx{}, nil).
// 		MaxTimes(1)
// 	c := Cursor{
// 		DB: s,
// 	}
// 	require.Error(t, c.AddBatch(parent, mockBatch))
// }
