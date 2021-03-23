// Package postgres provide faktory method for postgres db.Option
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/adipurnama/go-toolkit/db"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq" // use wrapped postgres driver
)

const intervalKeepAlive = 5 * time.Second

// NewPostgresDatabase - create & validate postgres connection given certain db.Option
// the caller have the responsibility to close the *sqlx.DB when succeed.
func NewPostgresDatabase(opt *db.Option) (*sql.DB, error) {
	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(opt.Username, opt.Password),
		Host:   fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Path:   opt.DatabaseName,
	}
	q := connURL.Query()
	q.Add("sslmode", "disable")
	connURL.RawQuery = q.Encode()

	db, err := apmsql.Open("postgres", connURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "postgres: failed to open connection")
	}

	db.SetMaxIdleConns(opt.ConnectionOption.MaxIdle)
	db.SetConnMaxLifetime(opt.ConnectionOption.MaxLifetime)
	db.SetMaxOpenConns(opt.ConnectionOption.MaxOpen)

	ctx, cancel := context.WithTimeout(context.Background(), opt.ConnectionOption.ConnectTimeout)
	defer cancel()

	_ = db.QueryRowContext(ctx, "SELECT 1")

	log.Println("successfully connected to postgres", connURL.Host)

	go doKeepAliveConnection(db, opt.DatabaseName, intervalKeepAlive)

	return db, nil
}

func doKeepAliveConnection(db *sql.DB, dbName string, interval time.Duration) {
	for {
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
			log.Printf("ERROR db.doKeepAliveConnection counter=%d db_name=%s stats=%v\n", i, dbName, db.Stats())
		}

		_ = rows.Close()

		time.Sleep(interval)
	}
}
