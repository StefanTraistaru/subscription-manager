package db

import (
	"fmt"
	"log"
	"os"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.in/mgo.v2"
)

func GetDBCollection() (*mongo.Collection, error) {
	// client, err := mongo.Connect(context.TODO(), "mongodb://mongo-credentials:27017")
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// client, err := mongo.Connect(context.TODO(), clientOptions)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo-credentials:27018"))
	if err != nil {
		fmt.Println("eroare db 1")
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("eroare db 2")
		return nil, err
	}
	collection := client.Database("GoLogin").Collection("users")
	return collection, nil
}

func GetDBCollection2() (*mgo.Collection, *mgo.Session, error) {
	// Connect to mongo
	session, err := mgo.Dial("mongo-credentials:27017")
	if err != nil {
		fmt.Println("db eroare")
		log.Fatalln(err)
		log.Fatalln("mongo err")
		os.Exit(1)
	}
	// defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// Get subscriptions collection
	collection := session.DB("app").C("users")

	return collection, session, nil
}