// Package postgres provide faktory method for postgres db.Option
package postgres

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/adipurnama/go-toolkit/db"
	"github.com/jmoiron/sqlx"
)

const intervalKeepAlive = 5 * time.Second

// NewPostgresDatabase - create & validate postgres connection given certain db.Option
// the caller have the responsibility to close the *sqlx.DB when succeed.
func NewPostgresDatabase(opt *db.Option) (*sqlx.DB, error) {
	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(opt.Username, opt.Password),
		Host:   fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Path:   opt.DatabaseName,
	}
	q := connURL.Query()
	q.Add("sslmode", "disable")
	connURL.RawQuery = q.Encode()

	db, err := sqlx.Open("postgres", connURL.String())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(opt.ConnectionOption.MaxIdle)
	db.SetConnMaxLifetime(opt.ConnectionOption.MaxLifetime)
	db.SetMaxOpenConns(opt.ConnectionOption.MaxOpen)

	_ = db.QueryRow("SELECT 1")

	go doKeepAliveConnection(db, opt.DatabaseName, intervalKeepAlive)

	return db, nil
}

func doKeepAliveConnection(db *sqlx.DB, dbName string, interval time.Duration) {
	for {
		rows, err := db.Query("SELECT 1")
		if err != nil {
			log.Printf("db.doKeepAliveConnection conn=postgres error=%s db_name=%s\n", err, dbName)
			return
		}

		if rows.Next() {
			var i int

			_ = rows.Scan(&i)
			log.Printf("db.doKeepAliveConnection counter=%d db_name=%s stats=%v\n", i, dbName, db.Stats())
		}

		_ = rows.Close()

		time.Sleep(interval)
	}
}
