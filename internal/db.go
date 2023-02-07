package metrics

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbInterface interface {
	InitDb() error
	Add(*JSONMetrics) error
	Get(string) (*JSONMetrics, error)
	Close()
	Ping()
}

type Cursor struct {
	DbInterface
	Db      *sql.DB
	Context context.Context
	IsValid bool
}

func NewCursor(link, adaptor string) (*Cursor, error) {
	db, err := sql.Open(adaptor, link)
	if err != nil {
		ErrorLog.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return &Cursor{
		Db:      db,
		Context: context.Background(),
		IsValid: true,
	}, nil
}

func (c *Cursor) Close() {
	c.Db.Close()
}

func (c *Cursor) Ping() bool {
	ctx, cancel := context.WithTimeout(c.Context, 1*time.Second)
	defer cancel()
	if err := c.Db.PingContext(ctx); err != nil {
		ErrorLog.Printf("ping error, database unreachable?: %e", err)
		return false
	}
	return true
}

func (c *Cursor) InitDb() error {
	_, err := c.Db.Exec(CREATE_GAUGE_TABLE)
	if err != nil {
		ErrorLog.Printf("error creating gaugemetrics table %e", err)
		return err
	}
	InfoLog.Println("gaugemetrics table was created")
	_, err = c.Db.Exec(CREATE_COUNTER_TABLE)
	if err != nil {
		ErrorLog.Printf("error creating countermetrics table %e", err)
		return err
	}
	InfoLog.Println("countermetrics table was created")

	return nil
}

func (c *Cursor) Add(incomingMetrics *JSONMetrics) error {
	switch incomingMetrics.MType {
	case "gauge":
		if row := c.Db.QueryRowContext(c.Context, INSERT_INTO_GAUGE, incomingMetrics.ID, incomingMetrics.MType, incomingMetrics.Value); row.Err() != nil {
			ErrorLog.Printf("error adding gauge row %s to db: %e", incomingMetrics.ID, row.Err())
			return row.Err()
		}
	case "counter":
		if row := c.Db.QueryRowContext(c.Context, INSERT_INTO_COUNTER, incomingMetrics.ID, incomingMetrics.MType, incomingMetrics.Delta); row.Err() != nil {
			ErrorLog.Printf("error adding counter row %s to db: %e", incomingMetrics.ID, row.Err())
			return row.Err()
		}
	}
	InfoLog.Printf("added %s data to db...", incomingMetrics.ID)
	return nil
}

func (c *Cursor) Get(metricToFind *JSONMetrics) (*JSONMetrics, error) {
	foundMetric := &JSONMetrics{}
	var row *sql.Row
	switch metricToFind.MType {
	case "gauge":
		if row = c.Db.QueryRowContext(c.Context, SELECT_FROM_GAUGE, metricToFind.ID); row.Err() != nil {
			ErrorLog.Printf("error getting gauge row %s to db: %e", metricToFind.ID, row.Err())
			return nil, row.Err()
		}
		err := row.Scan(foundMetric.ID, foundMetric.MType, foundMetric.Value)
		if err != nil {
			ErrorLog.Printf("error scanning gauge %s: %e", metricToFind.ID, err)
			return nil, err
		}
	case "counter":
		if row = c.Db.QueryRowContext(c.Context, SELECT_FROM_COUNTER, metricToFind.ID); row.Err() != nil {
			ErrorLog.Printf("error getting counter row %s to db: %e", metricToFind.ID, row.Err())
			return nil, row.Err()
		}
		err := row.Scan(foundMetric.ID, foundMetric.MType, foundMetric.Delta)
		if err != nil {
			ErrorLog.Printf("error scanning counter %s: %e", metricToFind.ID, err)
			return nil, err
		}
	}
	return foundMetric, nil
}
