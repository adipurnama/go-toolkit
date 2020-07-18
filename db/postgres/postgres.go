package postgres

import (
	"fmt"

	"github.com/adipurnama/go-toolkit/db"
	"github.com/jmoiron/sqlx"
)

// NewPostgresDatabase - create & validate postgres connection given certain db.Option
// the caller have the responsibility to close the *sqlx.DB when succeed.
func NewPostgresDatabase(opt *db.Option) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", opt.Username, opt.Password, opt.Host, opt.Port, opt.DatabaseName))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
