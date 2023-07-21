package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // required import for pgx

	"github.com/nmramorov/gowatcher/internal/collector/metrics"
	"github.com/nmramorov/gowatcher/internal/log"
)

var (
	GAUGE   = "gauge"
	COUNTER = "counter"

	DBDefaultTimeout = time.Duration(100) * time.Millisecond
)

type DriverMethods interface {
	Close() error
	PingContext(ctx context.Context) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type DatabaseAccess interface {
	InitDb(context.Context) error
	Add(context.Context, *metrics.JSONMetrics) error
	Get(context.Context, *metrics.JSONMetrics) (*metrics.JSONMetrics, error)
	CloseConnection(context.Context) error
	Ping(context.Context) error
	AddBatch(context.Context, *[]metrics.JSONMetrics) error
	Flush(context.Context) error
}

type Cursor struct {
	DatabaseAccess
	DB      DriverMethods
	IsValid bool
	buffer  []*metrics.JSONMetrics
}

func NewCursor(parent context.Context, link, adaptor string) (*Cursor, error) {
	// Add InitDB here
	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
	defer cancel()

	db, err := sql.Open(adaptor, link)
	if err != nil {
		log.ErrorLog.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	cursor := &Cursor{
		DB:      db,
		IsValid: true,
		buffer:  make([]*metrics.JSONMetrics, 0, 100),
	}

	err = cursor.Ping(ctx)
	if err != nil {
		cursor.IsValid = false
	}
	return cursor, nil
}

func (c *Cursor) CloseConnection(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
	defer cancel()

	err := func(ctx context.Context, cursor *Cursor) error {
		select {
		case <-ctx.Done():
			err := cursor.DB.Close()
			if err != nil {
				log.ErrorLog.Printf("error closing db: %e", err)
				return err
			}
			return err
		default:
			err := cursor.DB.Close()
			if err != nil {
				log.ErrorLog.Printf("error closing db: %e", err)
				return err
			}
			return err
		}
	}(ctx, c)

	return err
}

func (c *Cursor) Ping(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
	defer cancel()

	if err := c.DB.PingContext(ctx); err != nil {
		log.ErrorLog.Printf("ping error, database unreachable?: %e", err)
		return err
	}
	return nil
}

func (c *Cursor) InitDB(parent context.Context) error {
	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
	defer cancel()

	_, err := c.DB.ExecContext(ctx, CreateGaugeTable)
	if err != nil {
		log.ErrorLog.Printf("error creating gaugemetrics table %e", err)
		return err
	}
	log.InfoLog.Println("gaugemetrics table was created")
	_, err = c.DB.ExecContext(ctx, CreateCounterTable)
	if err != nil {
		log.ErrorLog.Printf("error creating countermetrics table %e", err)
		return err
	}
	log.InfoLog.Println("countermetrics table was created")

	return nil
}

func add(parent context.Context, incomingMetrics *metrics.JSONMetrics, db interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}) error {
	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
	defer cancel()
	switch incomingMetrics.MType {
	case GAUGE:
		if _, err := db.ExecContext(
			ctx, InsertIntoGauge, incomingMetrics.ID, incomingMetrics.MType, incomingMetrics.Value); err != nil {
			log.ErrorLog.Printf("error adding gauge row %s to DB: %e", incomingMetrics.ID, err)
			return err
		}
	case COUNTER:
		if _, err := db.ExecContext(
			ctx, InsertIntoCounter, incomingMetrics.ID, incomingMetrics.MType, incomingMetrics.Delta); err != nil {
			log.ErrorLog.Printf("error adding counter row %s to db: %e", incomingMetrics.ID, err)
			return err
		}
	}
	log.InfoLog.Printf("added %s data to db...", incomingMetrics.ID)
	return nil
}

func (c *Cursor) Add(parent context.Context, incomingMetrics *metrics.JSONMetrics) error {
	return add(parent, incomingMetrics, c.DB)
}

func (c *Cursor) Get(parent context.Context, metricToFind *metrics.JSONMetrics) (*metrics.JSONMetrics, error) {
	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
	defer cancel()

	foundMetric := &metrics.JSONMetrics{}
	var row *sql.Row
	switch metricToFind.MType {
	case GAUGE:
		if row = c.DB.QueryRowContext(ctx, SelectFromGauge, metricToFind.ID); row.Err() != nil {
			log.ErrorLog.Printf("error getting gauge row %s to db: %e", metricToFind.ID, row.Err())
			return nil, row.Err()
		}
		err := row.Scan(foundMetric.ID, foundMetric.MType, foundMetric.Value)
		if err != nil {
			log.ErrorLog.Printf("error scanning gauge %s: %e", metricToFind.ID, err)
			return nil, err
		}
	case COUNTER:
		if row = c.DB.QueryRowContext(ctx, SelectFromCounter, metricToFind.ID); row.Err() != nil {
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

func (c *Cursor) AddBatch(parent context.Context, metrics []*metrics.JSONMetrics) error {
	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
	defer cancel()

	c.buffer = append(c.buffer, metrics...)
	if cap(c.buffer) == len(c.buffer) {
		err := c.Flush(ctx)
		if err != nil {
			log.ErrorLog.Printf("cannot add record to the database")
			return err
		}
		c.buffer = c.buffer[:0]
	}
	return nil
}

func (c *Cursor) Flush(parent context.Context) error {
	return flush(parent, c.buffer, c.DB)
}

func flush(parent context.Context, buffer []*metrics.JSONMetrics, db interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}) error {
	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
	defer cancel()

	// проверим на всякий случай
	if db == nil {
		log.ErrorLog.Printf("You haven`t opened the database connection")
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmtGauge, err := tx.PrepareContext(ctx, InsertIntoGauge)
	if err != nil {
		return err
	}
	stmtCounter, err := tx.PrepareContext(ctx, InsertIntoCounter)
	if err != nil {
		return err
	}
	defer func() {
		err = stmtGauge.Close()
		if err != nil {
			log.ErrorLog.Printf("error closing Gauge statement: %e", err)
		}
	}()
	defer func() {
		err = stmtCounter.Close()
		if err != nil {
			log.ErrorLog.Printf("error closing Counter statement: %e", err)
		}
	}()

	for _, v := range buffer {
		switch v.MType {
		case GAUGE:
			if _, err = stmtGauge.ExecContext(ctx, v.ID, v.MType, v.Value); err != nil {
				if err = tx.Rollback(); err != nil {
					log.ErrorLog.Printf("update drivers: unable to rollback: %v", err)
				}
				return err
			}
		case COUNTER:
			if _, err = stmtCounter.ExecContext(ctx, v.ID, v.MType, v.Delta); err != nil {
				if err = tx.Rollback(); err != nil {
					log.ErrorLog.Printf("update drivers: unable to rollback: %v", err)
				}
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.ErrorLog.Printf("update drivers: unable to commit: %v", err)
		return err
	}
	return nil
}

// func (c *Cursor) Flush(parent context.Context) error {
// 	ctx, cancel := context.WithTimeout(parent, DBDefaultTimeout)
// 	defer cancel()

// 	// проверим на всякий случай
// 	if c.DB == nil {
// 		log.ErrorLog.Printf("You haven`t opened the database connection")
// 	}
// 	tx, err := c.DB.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}

// 	stmtGauge, err := tx.PrepareContext(ctx, InsertIntoGauge)
// 	if err != nil {
// 		return err
// 	}
// 	stmtCounter, err := tx.PrepareContext(ctx, InsertIntoCounter)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		err = stmtGauge.Close()
// 		if err != nil {
// 			log.ErrorLog.Printf("error closing Gauge statement: %e", err)
// 		}
// 	}()
// 	defer func() {
// 		err = stmtCounter.Close()
// 		if err != nil {
// 			log.ErrorLog.Printf("error closing Counter statement: %e", err)
// 		}
// 	}()

// 	for _, v := range c.buffer {
// 		switch v.MType {
// 		case GAUGE:
// 			if _, err = stmtGauge.ExecContext(ctx, v.ID, v.MType, v.Value); err != nil {
// 				if err = tx.Rollback(); err != nil {
// 					log.ErrorLog.Printf("update drivers: unable to rollback: %v", err)
// 				}
// 				return err
// 			}
// 		case COUNTER:
// 			if _, err = stmtCounter.ExecContext(ctx, v.ID, v.MType, v.Delta); err != nil {
// 				if err = tx.Rollback(); err != nil {
// 					log.ErrorLog.Printf("update drivers: unable to rollback: %v", err)
// 				}
// 				return err
// 			}
// 		}
// 	}

// 	if err := tx.Commit(); err != nil {
// 		log.ErrorLog.Printf("update drivers: unable to commit: %v", err)
// 		return err
// 	}

// 	c.buffer = c.buffer[:0]
// 	return nil
// }
