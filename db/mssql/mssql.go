package mssql

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb" // mssql driver
	"github.com/pkg/errors"

	"github.com/adipurnama/go-toolkit/db"
)

// NewMsSQLDatabase - create & validate mssql connection given certain db.Option
// the caller have the responsibility to close the *sql.DB when succeed.
func NewMsSQLDatabase(opt *db.Option) (*sql.DB, error) {
	connURL := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(opt.Username, opt.Password),
		Host:   fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		// Path:   opt.DatabaseName,
	}
	q := connURL.Query()
	q.Add("encrypt", "disable")
	q.Add("dial timeout", "10")
	q.Add("database", opt.DatabaseName)
	connURL.RawQuery = q.Encode()

	db, err := sql.Open("mssql", connURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "mssql: failed to open connection")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "mssql: error pinging database")
	}

	return db, nil
}
