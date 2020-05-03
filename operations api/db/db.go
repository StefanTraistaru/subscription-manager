package db

import (
    // "fmt"
    "log"
    "os"

    "gopkg.in/mgo.v2"
)

func GetDBCollection() (*mgo.Collection, *mgo.Session, error) {
    // Connect to mongo
    session, err := mgo.Dial("mongo:27017")
    if err != nil {
        log.Fatalln(err)
        log.Fatalln("mongo err")
        os.Exit(1)
    }
    session.SetMode(mgo.Monotonic, true)

    // Get users collection
    collection := session.DB("app").C("subscriptions")

    return collection, session, nil
}