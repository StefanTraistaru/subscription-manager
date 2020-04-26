package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "./model"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
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
    ID          bson.ObjectId `db:"id" json:"id,omitempty" bson:"_id"`
    Name        string `json:"name"`
    Price       string `json:"price"`
    Details     string `json:"details"`
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
    router.HandleFunc("/user", createUser).Methods("POST")
    router.HandleFunc("/subscriptions", createSubscription).Methods("POST")
    router.HandleFunc("/subscriptions", getSubscriptions).Methods("GET")

    http.ListenAndServe(":5000", cors.AllowAll().Handler(router))
    log.Println("Listening on port 5000...")
}


func createUser(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request create user")

    // Read body
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responseError2(w, "Cannot read request body",err, http.StatusBadRequest)
        return
    }

    var user model.User
    err = json.Unmarshal(data, &user)
    if err != nil {
        responseError2(w, "Cannot unmarshall body", err, http.StatusBadRequest)
        return
    }

    var newUser model.DBUser
    newUser.Username = user.Username
    newUser.FirstName = user.FirstName
    newUser.LastName = user.LastName

    err = subscriptions.Insert(newUser)
    if err != nil {
        responseError2(w, "Cannot insert new user", err, http.StatusInternalServerError)
        return
    }

    responseJSON2(w, "User created successfully", "User created successfully")
}


func createSubscription(w http.ResponseWriter, r *http.Request) {
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
    // asd := bson.NewObjectID()
    subscription.ID = bson.NewObjectId()
    // asd := bson.NewObjectId()
    // log.Println(asd)
    // Insert new subscription
    err = subscriptions.Insert(subscription)
    if err != nil {
        responseError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    responseJSON(w, subscription)
}

func getSubscriptions(w http.ResponseWriter, r *http.Request) {
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


func responseError2(w http.ResponseWriter, logMessage string, err error, code int) {
    log.Println("Error: " + logMessage)
    log.Println(err.Error())

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    var res model.ResponseResult
    res.Error = err.Error()
    json.NewEncoder(w).Encode(res)
}

func responseJSON2(w http.ResponseWriter, logMessage string, result string) {
    log.Println(logMessage)
    w.Header().Set("Content-Type", "application/json")
    var res model.ResponseResult
    res.Result = result
    json.NewEncoder(w).Encode(res)
}