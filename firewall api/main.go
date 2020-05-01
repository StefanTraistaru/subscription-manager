package main

import (
    "log"
    "net/http"
    "./controller"
    "github.com/rs/cors"
    "github.com/gorilla/mux"
)

func main() {
    // Set up routes
    router := mux.NewRouter()

    // Login api
    router.HandleFunc("/register", controller.Register).Methods("POST")
    router.HandleFunc("/login", controller.Login).Methods("POST")
    router.HandleFunc("/delete/{username}", controller.Delete).Methods("DELETE")
    // Operations api
    router.HandleFunc("/subscriptions/{username}", controller.CreateSubscription).Methods("POST")
    router.HandleFunc("/subscriptions/{username}", controller.GetSubscriptions).Methods("GET")
    router.HandleFunc("/subscription/{username}/{subName}", controller.GetSubscription).Methods("GET")

    http.ListenAndServe(":5003", cors.AllowAll().Handler(router))
    log.Println("Listening on port 5003...")
}

