package mongokit

import (
	"context"
	"fmt"
	"net/url"

	"github.com/adipurnama/go-toolkit/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDBClient returns new mongodb client using the go mongo-driver
func NewMongoDBClient(opt *db.Option) (*mongo.Client, error) {
	connURL := &url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(opt.Username, opt.Password),
		Host:   fmt.Sprintf("%s:%d", opt.Host, opt.Port),
	}
	q := connURL.Query()
	q.Add("authSource", opt.DatabaseName)
	connURL.RawQuery = q.Encode()

	clientOptions := options.Client()
	clientOptions.ApplyURI(connURL.String())
	clientOptions.SetConnectTimeout(opt.ConnectionOption.ConnectTimeout)
	clientOptions.SetMaxConnIdleTime(opt.ConnectionOption.MaxLifetime)
	clientOptions.SetMaxPoolSize(uint64(opt.ConnectionOption.MaxOpen))

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), opt.ConnectionOption.ConnectTimeout)
	defer cancel()

	err = client.Connect(ctx)

	return client, err
}
