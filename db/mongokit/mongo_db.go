package mongokit

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/adipurnama/go-toolkit/db"
	pmongo "github.com/pinpoint-apm/pinpoint-go-agent/plugin/mongodriver"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDBClient returns new mongodb client using the go mongo-driver.
func NewMongoDBClient(opt *db.Option, authDB string) (*mongo.Database, error) {
	connURL := &url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(opt.Username, opt.Password),
		Host:   fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Path:   "/",
	}
	q := connURL.Query()
	q.Add("authSource", authDB)
	connURL.RawQuery = q.Encode()

	clientOptions := options.Client()
	clientOptions.ApplyURI(connURL.String())
	clientOptions.SetConnectTimeout(opt.ConnectionOption.ConnectTimeout)
	clientOptions.SetMaxConnIdleTime(opt.ConnectionOption.MaxLifetime)
	clientOptions.SetMaxPoolSize(uint64(opt.ConnectionOption.MaxOpen))
	clientOptions.SetMonitor(pmongo.NewMonitor())

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "mongokit: cannot create new mongo client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), opt.ConnectionOption.ConnectTimeout)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "mongokit: cannot connect to client")
	}

	log.Println("successfully connected to mongo", connURL.Host)

	return client.Database(opt.DatabaseName), nil
}
