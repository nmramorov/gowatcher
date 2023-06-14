package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/log"
)

type DbInterface interface {
	InitDb() error
	Add(*metrics.JSONMetrics) error
	Get(string) (*metrics.JSONMetrics, error)
	Close()
	Ping()
	UpdateBatch()
}

type Cursor struct {
	DbInterface
	Db      *sql.DB
	Context context.Context
	IsValid bool
	buffer  []*metrics.JSONMetrics
}

func NewCursor(link, adaptor string) (*Cursor, error) {
	db, err := sql.Open(adaptor, link)
	if err != nil {
		log.ErrorLog.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	new := &Cursor{
		Db:      db,
		Context: context.Background(),
		IsValid: true,
		buffer:  make([]*metrics.JSONMetrics, 0, 100),
	}
	valid := new.Ping()
	if !valid {
		new.IsValid = false
	}
	return new, nil
}

func (c *Cursor) Close() {
	c.Db.Close()
}

func (c *Cursor) Ping() bool {
	ctx, cancel := context.WithTimeout(c.Context, 1*time.Second)
	defer cancel()
	if err := c.Db.PingContext(ctx); err != nil {
		log.ErrorLog.Printf("ping error, database unreachable?: %e", err)
		return false
	}
	return true
}

func (c *Cursor) InitDb() error {
	_, err := c.Db.Exec(CreateGaugeTable)
	if err != nil {
		log.ErrorLog.Printf("error creating gaugemetrics table %e", err)
		return err
	}
	log.InfoLog.Println("gaugemetrics table was created")
	_, err = c.Db.Exec(CreateCounterTable)
	if err != nil {
		log.ErrorLog.Printf("error creating countermetrics table %e", err)
		return err
	}
	log.InfoLog.Println("countermetrics table was created")

	return nil
}

func (c *Cursor) Add(incomingMetrics *metrics.JSONMetrics) error {
	switch incomingMetrics.MType {
	case "gauge":
		if _, err := c.Db.ExecContext(c.Context, InsertIntoGauge, incomingMetrics.ID, incomingMetrics.MType, incomingMetrics.Value); err != nil {
			log.ErrorLog.Printf("error adding gauge row %s to db: %e", incomingMetrics.ID, err)
			return err
		}
	case "counter":
		if _, err := c.Db.ExecContext(c.Context, InsertIntoCounter, incomingMetrics.ID, incomingMetrics.MType, incomingMetrics.Delta); err != nil {
			log.ErrorLog.Printf("error adding counter row %s to db: %e", incomingMetrics.ID, err)
			return err
		}
	}
	log.InfoLog.Printf("added %s data to db...", incomingMetrics.ID)
	return nil
}

func (c *Cursor) Get(metricToFind *metrics.JSONMetrics) (*metrics.JSONMetrics, error) {
	foundMetric := &metrics.JSONMetrics{}
	var row *sql.Row
	switch metricToFind.MType {
	case "gauge":
		if row = c.Db.QueryRowContext(c.Context, SelectFromGauge, metricToFind.ID); row.Err() != nil {
			log.ErrorLog.Printf("error getting gauge row %s to db: %e", metricToFind.ID, row.Err())
			return nil, row.Err()
		}
		err := row.Scan(foundMetric.ID, foundMetric.MType, foundMetric.Value)
		if err != nil {
			log.ErrorLog.Printf("error scanning gauge %s: %e", metricToFind.ID, err)
			return nil, err
		}
	case "counter":
		if row = c.Db.QueryRowContext(c.Context, SelectFromCounter, metricToFind.ID); row.Err() != nil {
			log.ErrorLog.Printf("error getting counter row %s to db: %e", metricToFind.ID, row.Err())
			return nil, row.Err()
		}
		err := row.Scan(foundMetric.ID, foundMetric.MType, foundMetric.Delta)
		if err != nil {
			log.ErrorLog.Printf("error scanning counter %s: %e", metricToFind.ID, err)
			return nil, err
		}
	}
	return foundMetric, nil
}

func (c *Cursor) AddBatch(metrics []*metrics.JSONMetrics) error {
	c.buffer = append(c.buffer, metrics...)
	if cap(c.buffer) == len(c.buffer) {
		err := c.Flush()
		if err != nil {
			log.ErrorLog.Printf("cannot add record to the database")
			return err
		}
	}
	return nil
}

func (c *Cursor) Flush() error {
	// проверим на всякий случай
	if c.Db == nil {
		log.ErrorLog.Printf("You haven`t opened the database connection")
	}
	tx, err := c.Db.Begin()
	if err != nil {
		return err
	}

	stmtGauge, err := tx.Prepare(InsertIntoGauge)
	if err != nil {
		return err
	}
	stmtCounter, err := tx.Prepare(InsertIntoCounter)
	if err != nil {
		return err
	}
	defer stmtGauge.Close()
	defer stmtCounter.Close()

	for _, v := range c.buffer {
		switch v.MType {
		case "gauge":
			if _, err = stmtGauge.Exec(v.ID, v.MType, v.Value); err != nil {
				if err = tx.Rollback(); err != nil {
					log.ErrorLog.Fatalf("update drivers: unable to rollback: %v", err)
				}
				return err
			}
		case "counter":
			if _, err = stmtCounter.Exec(v.ID, v.MType, v.Delta); err != nil {
				if err = tx.Rollback(); err != nil {
					log.ErrorLog.Fatalf("update drivers: unable to rollback: %v", err)
				}
				return err
			}
		}

	}

	if err := tx.Commit(); err != nil {
		log.ErrorLog.Fatalf("update drivers: unable to commit: %v", err)
		return err
	}

	c.buffer = c.buffer[:0]
	return nil
}
