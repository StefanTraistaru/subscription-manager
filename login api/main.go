package main

import (
    "./controller"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/register", controller.RegisterHandler).Methods("POST")
    r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
    r.HandleFunc("/profile", controller.ProfileHandler).Methods("GET")

    log.Fatal(http.ListenAndServe(":5002", r))
}