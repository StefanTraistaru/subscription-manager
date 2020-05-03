package controller

import (
    "encoding/json"
    "log"
    "io/ioutil"
    "net/http"
    "errors"

    "../db"
    "../model"

    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/mux"
)



func CreateUser(w http.ResponseWriter, r *http.Request) {
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

    // Get connection to DB
    collection, session, err := db.GetDBCollection()
    if err != nil {
        responseError(w, "Cannot get connection to DB", err, http.StatusInternalServerError)
        return
    }
    defer session.Close()

    // Create and add user to DB
    newUser := model.DBUser {
        Username: user.Username,
        FirstName: user.FirstName,
        LastName: user.LastName,
    }

    err = collection.Insert(newUser)
    if err != nil {
        responseError(w, "Cannot insert new user", err, http.StatusInternalServerError)
        return
    }

    // Create and send response
    var response model.ResponseResult
    response.Result = "User created successfully"
    responseJSON(w, "Operation successful", response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: delete")

    // Read body
    params := mux.Vars(r)
    id := params["username"]

    // Get connection to DB
    collection, session, err := db.GetDBCollection()
    if err != nil {
        responseError(w, "Cannot get connection to DB", err, http.StatusInternalServerError)
        return
    }
    defer session.Close()

    // Query DB for user
    err = collection.Remove(bson.M{"username": id})
    if err != nil {
        if err.Error() == "not found" {
            responseJSON(w, "Invalid user", "Invalid user")
            return
        }
        responseError(w, "Query DB", err, http.StatusInternalServerError)
        return
    }

    // Create and send response
    var response model.ResponseResult
    response.Result = "Deletion successful"
    responseJSON(w, "Deletion successful", response)
}

func CreateSubscription(w http.ResponseWriter, r *http.Request) {
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

    // Get connection to DB
    collection, session, err := db.GetDBCollection()
    if err != nil {
        responseError(w, "Cannot get connection to DB", err, http.StatusInternalServerError)
        return
    }
    defer session.Close()

    // Find user in DB
    var result model.DBUser
    err = collection.Find(bson.M{"username": id}).One(&result)
    if err != nil {
        responseError(w, "Cannot query DB", err, http.StatusInternalServerError)
        return
    }

    // Create and add subscription to this user
    subscription.ID = bson.NewObjectId()
    result.Subscriptions = append(result.Subscriptions, subscription)
    err = collection.Update(bson.M{"username": id}, result)
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

func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request get all subscriptions")

    // Read body
    params := mux.Vars(r)
    username := params["username"]

    // Get connection to DB
    collection, session, err := db.GetDBCollection()
    if err != nil {
        responseError(w, "Cannot get connection to DB", err, http.StatusInternalServerError)
        return
    }
    defer session.Close()

    // Find user in DB
    var result model.DBUser
    err = collection.Find(bson.M{"username": username}).One(&result)
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

func GetSubscription(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request get one subscription")

    // Read body
    params := mux.Vars(r)
    username := params["username"]
    subscriptionName := params["subName"]

    // Get connection to DB
    collection, session, err := db.GetDBCollection()
    if err != nil {
        responseError(w, "Cannot get connection to DB", err, http.StatusInternalServerError)
        return
    }
    defer session.Close()

    // Find user in DB
    var result model.DBUser
    err = collection.Find(bson.M{"username": username}).One(&result)
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