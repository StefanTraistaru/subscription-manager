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
    router.HandleFunc("/user", controller.CreateUser).Methods("POST")
    router.HandleFunc("/user/delete/{username}", controller.DeleteUser).Methods("POST")
    router.HandleFunc("/subscriptions/{username}", controller.CreateSubscription).Methods("POST")
    router.HandleFunc("/subscriptions/{username}", controller.GetSubscriptions).Methods("GET")
    router.HandleFunc("/subscription/{username}/{subName}", controller.GetSubscription).Methods("GET")
    // TODO:
    // router.HandleFunc("/subscription/{username}/{subName}", updateSubscription).Methods("PUT")
    // router.HandleFunc("/subscription/{username}/{subName}", deleteSubscription).Methods("DELETE")

    http.ListenAndServe(":5000", cors.AllowAll().Handler(router))
    log.Println("Listening on port 5000...")
}
