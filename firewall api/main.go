package main

import (
    // "encoding/json"
    // "io/ioutil"
    "log"
    "net/http"
    // "bytes"
    "./controller"

    // "gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
    "github.com/rs/cors"
    "github.com/gorilla/mux"
)

type User struct {
    Username  string `json:"username"`
    FirstName string `json:"firstname"`
    LastName  string `json:"lastname"`
    Password  string `json:"password"`
    Token     string `json:"token"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

func main() {
    // Set up routes
    router := mux.NewRouter()
    router.HandleFunc("/register", controller.Register).Methods("POST")
    router.HandleFunc("/login", controller.Login).Methods("POST")
    router.HandleFunc("/subscriptions/{username}", controller.CreateSubscription).Methods("POST")
    router.HandleFunc("/subscriptions/{username}", controller.GetSubscriptions).Methods("GET")
    router.HandleFunc("/subscription/{username}/{subName}", controller.GetSubscription).Methods("GET")

    http.ListenAndServe(":5003", cors.AllowAll().Handler(router))
    log.Println("Listening on port 5003...")
}

