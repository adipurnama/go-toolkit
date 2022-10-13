package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/pkg/errors"
	_ "github.com/sijms/go-ora/v2" // use wrapped oracle driver

	"github.com/adipurnama/go-toolkit/db"
)

// NewOracleDatabase - create & validate postgres connection given certain db.Option
// the caller have the responsibility to close the *sqlx.DB when succeed.
func NewOracleDatabase(opt *db.Option, opts ...db.Options) (*sql.DB, error) {
	for _, o := range opts {
		o(opt)
	}

	connURL := &url.URL{
		Scheme: "oracle",
		User:   url.UserPassword(opt.Username, opt.Password),
		Host:   fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Path:   opt.DatabaseName,
	}
	q := connURL.Query()
	q.Add("sslmode", "disable")
	connURL.RawQuery = q.Encode()

	oraDB, err := sql.Open("oracle", connURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "oracle: failed to open connection")
	}

	oraDB.SetMaxIdleConns(opt.MaxIdle)
	oraDB.SetConnMaxLifetime(opt.MaxLifetime)
	oraDB.SetMaxOpenConns(opt.MaxOpen)

	ctx, cancel := context.WithTimeout(context.Background(), opt.ConnectTimeout)
	defer cancel()

	err = oraDB.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "oracle: failed to ping")
	}

	_ = oraDB.QueryRowContext(ctx, "SELECT 1 FROM DUAL")

	log.Println("successfully connected to oracle", connURL.Host)

	if opt.AppContext != nil {
		go doKeepAliveConnection(opt.AppContext, oraDB, opt.DatabaseName, opt.KeepAliveCheckInterval)
	}

	return oraDB, nil
}

func doKeepAliveConnection(ctx context.Context, db *sql.DB, dbName string, interval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			rows, err := db.Query("SELECT 1 FROM DUAL")
			if err != nil {
				log.Printf("ERROR db.doKeepAliveConnection conn=oracle error=%s db_name=%s\n", err, dbName)
				return
			}

			if rows.Err() != nil {
				log.Printf("ERROR db.doKeepAliveConnection conn=oracle error=%s db_name=%s\n", rows.Err(), dbName)
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
