package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func isInReleaseMode() bool {
	dbArg := ""
	if len(os.Args) != 1 {
		dbArg = os.Args[1]
	}
	return dbArg == "--release"
}

func DbInstace() *mongo.Client {

	URL, ok := os.LookupEnv("MONGO_URL")
	if !ok {
		panic("URL not found")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(URL))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Db Connected")
	return client
}

func OpenCollection(client *mongo.Client, dbName, collenctionName string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collenctionName)
	return collection
}
