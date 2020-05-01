package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "errors"

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
    router.HandleFunc("/subscriptions/{username}", createSubscription).Methods("POST")
    router.HandleFunc("/subscriptions/{username}", getSubscriptions).Methods("GET")
    router.HandleFunc("/subscription/{username}/{subName}", getSubscription).Methods("GET")

    http.ListenAndServe(":5000", cors.AllowAll().Handler(router))
    log.Println("Listening on port 5000...")
}


func createUser(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request create user")

    // Read body
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responseError(w, "Cannot read request body",err, http.StatusBadRequest)
        return
    }

    // Unmarshal body
    var user model.User
    err = json.Unmarshal(data, &user)
    if err != nil {
        responseError(w, "Cannot unmarshall body", err, http.StatusBadRequest)
        return
    }

    // Create and add user to DB
    newUser := model.DBUser {
        Username: user.Username,
        FirstName: user.FirstName,
        LastName: user.LastName,
    }
    err = subscriptions.Insert(newUser)
    if err != nil {
        responseError(w, "Cannot insert new user", err, http.StatusInternalServerError)
        return
    }

    // Create and send response
    var response model.ResponseResult
    response.Result = "User created successfully"
    responseJSON(w, "Operation successful", response)
}


func createSubscription(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request create subscription")

    // Read body
    params := mux.Vars(r)
    id := params["username"]
    log.Println("username: " + id)
    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        responseError(w, "Cannot read request body", err, http.StatusBadRequest)
        return
    }

    // Unmarshal body
    var subscription model.Subscription
    err = json.Unmarshal(data, &subscription)
    if err != nil {
        responseError(w, "Cannot unmarshal body", err, http.StatusBadRequest)
        return
    }

    // Find user in DB
    var result model.DBUser
    err = subscriptions.Find(bson.M{"username": id}).One(&result)
    if err != nil {
        responseError(w, "Cannot query DB", err, http.StatusInternalServerError)
        return
    }

    // Create and add subscription to this user
    subscription.ID = bson.NewObjectId()
    result.Subscriptions = append(result.Subscriptions, subscription)
    err = subscriptions.Update(bson.M{"username": id}, result)
    if err != nil {
        responseError(w, "Cannot add subscription in DB", err, http.StatusInternalServerError)
        return
    }

    // Create and send response
    var response model.ResponseResult
    var responseData []model.Subscription
    responseData = append(responseData, subscription)
    response.Data = responseData
    response.Result = "Successful"
    responseJSON(w, "Operation successful", response)
}


func getSubscriptions(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request get all subscriptions")

    // Read body
    params := mux.Vars(r)
    username := params["username"]

    // Find user in DB
    var result model.DBUser
    err := subscriptions.Find(bson.M{"username": username}).One(&result)
    if err != nil {
        responseError(w, "User not found", err, http.StatusBadRequest)
        return
    }

    // Create and send response
    var response model.ResponseResult
    response.Result = "Successful"
    response.Data = result.Subscriptions
    responseJSON(w, "Operation successful", response)
}

func getSubscription(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request get one subscription")

    // Read body
    params := mux.Vars(r)
    username := params["username"]
    subscriptionName := params["subName"]

    // Find user in DB
    var result model.DBUser
    err := subscriptions.Find(bson.M{"username": username}).One(&result)
    if err != nil {
        responseError(w, "User not found", err, http.StatusBadRequest)
        return
    }

    // Find subscription
    for _, sub := range result.Subscriptions {
        if sub.Name == subscriptionName {
            var response model.ResponseResult
            var responseData []model.Subscription
            responseData = append(responseData, sub)
            response.Data = responseData
            response.Result = "Successful"
            responseJSON(w, "Operation successful", response)
            return
        }
    }

    // TODO: Query the DB for subscription directly
    // Testing ---------------------------
    // var result2 model.Subscription
    // err = subscriptions.Find(bson.M{"username": username, "subscriptions.name": subscriptionName}).One(&result)
    // if err != nil {
    //     log.Println("Error query DB")
    //     log.Println(err.Error())
    //     return
    // }
    // fmt.Println(result2)
    // responseJSON(w, result2)
    // -----------------------------------

    responseError(w, "Subscription not found", errors.New("Subscription not found"), http.StatusBadRequest)
}

////////////////////////////////////////////////////
// Util functions
////////////////////////////////////////////////////

func responseJSON(w http.ResponseWriter, message string, data interface{}) {
    log.Println(message)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func responseError(w http.ResponseWriter, logMessage string, err error, code int) {
    log.Println("Error: " + logMessage)
    log.Println(err.Error())

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    var res model.ResponseResult
    res.Error = err.Error()
    json.NewEncoder(w).Encode(res)
}