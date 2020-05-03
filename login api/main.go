package main

import (
    "./controller"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/register", controller.Register).Methods("POST")
    r.HandleFunc("/login", controller.Login).Methods("POST")
    r.HandleFunc("/delete/{username}", controller.Delete).Methods("POST")

    log.Fatal(http.ListenAndServe(":5002", r))
}