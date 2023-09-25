package database

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	errorshandler "github.com/crawler-tokens/migration/utils"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDataBase() *mongo.Client {
	DB_CONNECT_URL, ok := viper.Get("MONGO_CONNECTION_URI").(string)

	if !ok {
		log.Fatalf("error while fetching db connect url from yaml file: %v", ok)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(DB_CONNECT_URL))

	if err != nil {
		e := errorshandler.GetErrors(err, http.StatusInternalServerError)
		log.Fatal(e.HandleError())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		e := errorshandler.GetErrors(err, http.StatusInternalServerError)
		log.Fatal(e.HandleError())
	}

	fmt.Println("DB connected successfully")
	ClientConnection = client
	return client
}

var ClientConnection *mongo.Client

func OpenCollection(collectionName string) *mongo.Collection {
	return ClientConnection.Database("KoinPark_migration").Collection(collectionName)
}
