package main

import (
	"context"

	"github.com/adipurnama/go-toolkit/db"
	"github.com/adipurnama/go-toolkit/db/mongokit"
	"github.com/adipurnama/go-toolkit/log"
)

type doc struct {
	Name string `bson:"name"`
	Desc string `bson:"desc"`
}

func main() {
	// docker run -it -v mongodata:/data/db -p 27017:27017 --name mongodb -d mongo
	opt, err := db.NewDatabaseOption("127.0.0.1", 27017, "root", "example", "sample_db", nil)
	if err != nil {
		log.Fatal("error found while building mongodb connection ", err)
	}

	client, err := mongokit.NewMongoDBClient(opt, "admin")
	if err != nil {
		log.Fatal("found error while connecting to mongo ", err)
	}

	res, err := client.Collection("my-collection").InsertOne(context.Background(), doc{
		Name: "my-name",
		Desc: "my description",
	})
	if err != nil {
		log.Fatal("found error inserting document ", err)
	}

	log.Println("inserted with ID: ", res.InsertedID)
}
