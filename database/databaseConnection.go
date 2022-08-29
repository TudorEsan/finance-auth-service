package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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

	envs, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("no .env")
	}
	var mongoUrl string
	if isInReleaseMode() {
		mongoUrl = envs["MONGO_REALEASE_URL"]
	} else {
		mongoUrl = envs["MONGO_URL"]
	}
	if mongoUrl == "" {
		panic("Credentials not found")
	}
	client, error := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
	if error != nil {
		log.Fatal(error)
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

var Client *mongo.Client = DbInstace()

func OpenCollection(client *mongo.Client, collenctionName string) *mongo.Collection {
	collection := client.Database("Cluster0").Collection(collenctionName)
	return collection
}
