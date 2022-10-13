// Package postgres provide faktory method for postgres db.Option
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/pinpoint-apm/pinpoint-go-agent/plugin/pgsql" // use wrapped postgres driver
	"github.com/pkg/errors"

	"github.com/adipurnama/go-toolkit/db"
)

// NewPostgresDatabase - create & validate postgres connection given certain db.Option
// the caller have the responsibility to close the *sqlx.DB when succeed.
func NewPostgresDatabase(opt *db.Option, opts ...db.Options) (*sql.DB, error) {
	for _, o := range opts {
		o(opt)
	}

	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(opt.Username, opt.Password),
		Host:   fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Path:   opt.DatabaseName,
	}
	q := connURL.Query()
	q.Add("sslmode", "disable")
	connURL.RawQuery = q.Encode()

	db, err := sql.Open("pq-pinpoint", connURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to open connection")
	}

	db.SetMaxIdleConns(opt.MaxIdle)
	db.SetConnMaxLifetime(opt.MaxLifetime)
	db.SetMaxOpenConns(opt.MaxOpen)

	ctx, cancel := context.WithTimeout(context.Background(), opt.ConnectTimeout)
	defer cancel()

	_ = db.QueryRowContext(ctx, "SELECT 1")

	log.Println("successfully connected to postgres", connURL.Host)

	if opt.AppContext != nil {
		go doKeepAliveConnection(opt.AppContext, db, opt.DatabaseName, opt.KeepAliveCheckInterval)
	}

	return db, nil
}

func doKeepAliveConnection(ctx context.Context, db *sql.DB, dbName string, interval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			rows, err := db.Query("SELECT 1")
			if err != nil {
				log.Printf("ERROR db.doKeepAliveConnection conn=postgres error=%s db_name=%s\n", err, dbName)
				return
			}

			if rows.Err() != nil {
				log.Printf("ERROR db.doKeepAliveConnection conn=postgres error=%s db_name=%s\n", rows.Err(), dbName)
				return
			}

			if rows.Next() {
				var i int

				_ = rows.Scan(&i)
				log.Printf("SUCCESS db.doKeepAliveConnection counter=%d db_name=%s stats=%+v\n", i, dbName, db.Stats())
			}

			_ = rows.Close()

			time.Sleep(interval)
		}
	}
}
