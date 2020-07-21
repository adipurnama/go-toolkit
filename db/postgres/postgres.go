package postgres

import (
	"fmt"
	"net/url"
	"time"

	"github.com/adipurnama/go-toolkit/db"
	"github.com/adipurnama/go-toolkit/log"
	"github.com/jmoiron/sqlx"
)

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

	_ = db.QueryRow("SELECT 1")

	go doKeepAliveConnection(db, 5*time.Second)

	return db, nil
}

func doKeepAliveConnection(db *sqlx.DB, interval time.Duration) {
	for {
		rows, err := db.Query("SELECT 1")
		if err != nil {
			log.Error(err, "db.doKeepAliveConnection", "conn", "postgres")
			return
		}

		if rows.Next() {
			var i int

			_ = rows.Scan(&i)
			log.Info("db.doKeepAliveConnection", "counter", i, "stats", db.Stats())
		}

		_ = rows.Close()

		time.Sleep(interval)
	}
}
