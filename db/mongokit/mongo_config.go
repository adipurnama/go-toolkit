// Package mongokit provides helper to interact with mongodb storage
package mongokit

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/adipurnama/go-toolkit/config"
	"github.com/adipurnama/go-toolkit/db"
)

/*
NewFromConfig returns mongo *mongo.Database instance from yaml config file

	given config file contents:

		db:
		  primary:
		    username: <username>
		    password: "<password>"
		    host: mymongo.host.com
		    port: 270017
		    schema: mymongo_db
		    auth-db: admin
		    conn:
			  max-idle: 20
			  max-lifetime: 10m
			  timeout: 5m
			  max-open: 100

	 then we can call using :

		 v := viper.New()
		 ... set v file configs, etc

		 db, err := log.NewFromConfig(v, "db.primary")
		 ...continue using db.
*/
func NewFromConfig(v config.KVStore, path string) (*mongo.Database, error) {
	connOpt := db.DefaultConnectionOption()

	if maxIdle := v.GetInt(fmt.Sprintf("%s.conn.max-idle", path)); maxIdle > 0 {
		connOpt.MaxIdle = v.GetInt(fmt.Sprintf("%s.conn.max-idle", path))
	}

	if maxOpen := v.GetInt(fmt.Sprintf("%s.conn.max-open", path)); maxOpen > 0 {
		connOpt.MaxOpen = maxOpen
	}

	if maxLifetime := v.GetDuration(fmt.Sprintf("%s.conn.max-lifetime", path)); maxLifetime > 0 {
		connOpt.MaxLifetime = maxLifetime
	}

	if connTimeout := v.GetDuration(fmt.Sprintf("%s.conn.timeout", path)); connTimeout > 0 {
		connOpt.ConnectTimeout = connTimeout
	}

	opt, err := db.NewDatabaseOption(
		v.GetString(fmt.Sprintf("%s.host", path)),
		v.GetInt(fmt.Sprintf("%s.port", path)),
		v.GetString(fmt.Sprintf("%s.username", path)),
		v.GetString(fmt.Sprintf("%s.password", path)),
		v.GetString(fmt.Sprintf("%s.schema", path)),
		connOpt,
	)
	if err != nil {
		return nil, err
	}

	authDB := v.GetString(fmt.Sprintf("%s.auth-db", path))

	return NewMongoDBClient(opt, authDB)
}
