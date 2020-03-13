package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"

	"gopkg.in/mgo.v2"
    "github.com/rs/cors"
	"github.com/gorilla/mux"
)

const (
    hosts      = "subscription-manager_mongodb_1:27017"
    database   = "db"
    username   = ""
    password   = ""
    collection = "jobs"
)

type Subscription struct {
    Name        string `json:"name"`
    Price       string `json:"price"`
    Description string `json:"description"`
    Date_d      string `json:"date_d"`
    Date_m      string `json:"date_m"`
    Date_y      string `json:"date_y"`
}

var subscriptions *mgo.Collection

func main() {
    // Connect to mongo
    session, err := mgo.Dial("mongo:27017")
    if err != nil {
        log.Fatalln(err)
        log.Fatalln("mongo err")
        os.Exit(1)
    }
    defer session.Close()
    session.SetMode(mgo.Monotonic, true)

    // Get subscriptions collection
    subscriptions = session.DB("app").C("subscriptions")

    // Set up routes
    router := mux.NewRouter()
    router.HandleFunc("/subscriptions", createSubscription).Methods("POST")
    router.HandleFunc("/subscriptions", getSubscriptions).Methods("GET")

    http.ListenAndServe(":5000", cors.AllowAll().Handler(router))
    log.Println("Listening on port 5000...")
}

func createSubscription(w http.ResponseWriter, r *http.Request) {
    // TODO: compare with supply chain solution method
    // Read body
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responseError(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Read post
    subscription := &Subscription{}
    err = json.Unmarshal(data, subscription)
    if err != nil {
        responseError(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Insert new subscription
    err = subscriptions.Insert(subscription)
    if err != nil {
        responseError(w, err.Error(), http.StatusInternalServerError)
		return
    }

    responseJSON(w, subscription)
}

func getSubscriptions(w http.ResponseWriter, r *http.Request) {
    // TODO: what is this?
    result := []Subscription{}
    err := subscriptions.Find(nil).Sort("-name").All(&result)
    if err != nil {
        responseError(w, err.Error(), http.StatusInternalServerError)
    } else {
        responseJSON(w, result)
    }
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// TODO: Should I use this instead?
// func initialiseMongo() (session *mgo.Session){

//     info := &mgo.DialInfo{
//         Addrs:    []string{hosts},
//         Timeout:  60 * time.Second,
//         Database: database,
//         Username: username,
//         Password: password,
//     }

//     session, err := mgo.DialWithInfo(info)
//     if err != nil {
//         panic(err)
//     }

//     return
// }


// ------------ OLD ------------
// func jobsGetHandler(w http.ResponseWriter, r *http.Request) {

//     col := mongoStore.session.DB(database).C(collection)

//     results := []Job{}
//     col.Find(bson.M{"title": bson.RegEx{"", ""}}).All(&results)
//     jsonString, err := json.Marshal(results)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Fprint(w, string(jsonString))

// }

// func jobsPostHandler(w http.ResponseWriter, r *http.Request) {

//     col := mongoStore.session.DB(database).C(collection)

//     //Retrieve body from http request
//     b, err := ioutil.ReadAll(r.Body)
//     defer r.Body.Close()
//     if err != nil {
//         panic(err)
//     }

//     //Save data into Job struct
//     var _job Job
//     err = json.Unmarshal(b, &_job)
//     if err != nil {
//         http.Error(w, err.Error(), 500)
//         return
//     }

//     //Insert job into MongoDB
//     err = col.Insert(_job)
//     if err != nil {
//         panic(err)
//     }

//     //Convert job struct into json
//     jsonString, err := json.Marshal(_job)
//     if err != nil {
//         http.Error(w, err.Error(), 500)
//         return
//     }

//     //Set content-type http header
//     w.Header().Set("content-type", "application/json")

//     //Send back data as response
//     w.Write(jsonString)
// }





// func handler(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprintf(w, "Hello %s!", r.URL.Path)
//     fmt.Println("RESTfulServ. on:8093, Controller:",r.URL.Path[1:])
// }

// func main() {

//     // Initialize router
//     router := mux.NewRouter()

//     // Define handler functions
//     router.HandleFunc("/", handler).Methods("POST")

//     fmt.Println("REST API")
//     log.Fatal(http.ListenAndServe(":5000", router))
// }
