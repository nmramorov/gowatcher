package metrics

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Cursor struct {
	Db      *sql.DB
	Context context.Context
}

func NewCursor(link string) *Cursor {
	db, err := sql.Open("pgx", link)
	if err != nil {
		ErrorLog.Printf("Unable to connect to database: %v\n", err)
	}
	return &Cursor{
		Db:      db,
		Context: context.Background(),
	}
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
