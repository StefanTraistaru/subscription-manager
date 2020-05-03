package main

import (
    "log"
    "net/http"
    "./controller"
    "github.com/rs/cors"
    "github.com/gorilla/mux"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type httpHandlerFunc func(http.ResponseWriter, *http.Request)

var counter = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Namespace: "myAPI",
        Name:      "requests_total",
        Help:      "Total number of requests.",
    }, []string{"endpoint"})

func main() {
    prometheus.MustRegister(counter)
    // Set up routes
    router := mux.NewRouter()

    // Login api
    router.HandleFunc("/register", LogPrometheus("Register", (controller.Register))).Methods("POST")
    router.HandleFunc("/login", LogPrometheus("Login", (controller.Login))).Methods("POST")
    router.HandleFunc("/delete/{username}", LogPrometheus("Delete user", (controller.Delete))).Methods("DELETE")
    // Operations api
    router.HandleFunc("/subscriptions/{username}", LogPrometheus("Create subscription", (controller.CreateSubscription))).Methods("POST")
    router.HandleFunc("/subscriptions/{username}", LogPrometheus("Get all subscriptions", (controller.GetSubscriptions))).Methods("GET")
    router.HandleFunc("/subscription/{username}/{subName}", LogPrometheus("Get one subscription", (controller.GetSubscription))).Methods("GET")
    // Prometheus metrics
    router.Handle("/metrics", promhttp.Handler())
    prometheus.Register(counter)

    http.ListenAndServe(":5003", cors.AllowAll().Handler(router))
    log.Println("Listening on port 5003...")
}

func LogPrometheus(request string, next httpHandlerFunc) httpHandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Increase prometheus statistics
        counter.WithLabelValues(request).Inc()
        // Call the next handler.
        next(w, r)
    }
}