package metrics

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbInterface interface {
	InitDb() error
	Update(*JSONMetrics) error
	Get(string) (*JSONMetrics, error)
	Close()
	Ping()
}

type Cursor struct {
	DbInterface
	Db      *sql.DB
	Context context.Context
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

func (c *Cursor) Update(incomingMetrics *JSONMetrics) error {
	return nil
}

func Get(metricName string) (*JSONMetrics, error) {
	return nil, nil
}
