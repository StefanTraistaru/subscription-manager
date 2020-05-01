package controller

import (
    "encoding/json"
    "log"
    "errors"
    // "fmt"
    // "../db"
    "bytes"
    "../model"
    "io/ioutil"
    "net/http"

    "github.com/gorilla/mux"
    // jwt "github.com/dgrijalva/jwt-go"
    // "go.mongodb.org/mongo-driver/bson"
    // "golang.org/x/crypto/bcrypt"
)

////////////////////////////////////////////////////
// Login api
////////////////////////////////////////////////////

func Register(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: register")

    // Read body of the request
    var user model.User
    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(body, &user)
    if err != nil {
        responseError(w, "Cannot unmarshal body", err, http.StatusBadRequest)
        return
    }

    // Sending request to login-api
    response, err := http.Post("http://login-api:5002/register", "application/json", bytes.NewReader(body))
    if err != nil {
        responseError(w, "Cannot send request to login-api", err, http.StatusInternalServerError)
        return
    }

    // Unmarshal response from login-api
    var res model.ResponseResult
    body, _ = ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &res)
    if err != nil {
        responseError(w, "Cannot unmarshal response from login-api", err, http.StatusInternalServerError)
        return
    }

    if res.Result == "Registration successful" {
        // Sending request to operations-api to create user
        user.Password = ""
        user.Token = ""
        jsonValue, _ := json.Marshal(user)
        response, err := http.Post("http://operations-api:5000/user", "application/json", bytes.NewBuffer(jsonValue))
        if err != nil {
            responseError(w, "Cannot send request to operations-api", err, http.StatusInternalServerError)
            return
        }

        // Unmarshall response from operations-api
        var res model.ResponseResult
        body, _ = ioutil.ReadAll(response.Body)
        err = json.Unmarshal(body, &res)
        if err != nil {
            responseError(w, "Cannot unmarshal response from login-api", err, http.StatusInternalServerError)
            return
        }

        if res.Result == "User created successfully" {
            responseJSON(w, "Registration successful", "Registration successful")
            return
        }

        responseError(w, "Cannot create user", nil, http.StatusInternalServerError)
        return
    }

    responseError(w, "Cannot register user", nil, http.StatusInternalServerError)
    return
}

func Login(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: login")

    // Read body of the request
    var user model.User
    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(body, &user)
    if err != nil {
        responseError(w, "Cannot unmarshal body", err, http.StatusBadRequest)
        return
    }

    // Sending request to login-api
    response, err := http.Post("http://login-api:5002/login", "application/json", bytes.NewReader(body))
    if err != nil {
        log.Println("Error in sending request to login-api")
        log.Println(err)
        return
    }

    // Unmarshall response from login-api
    var res model.ResponseResult
    body, _ = ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &res)
    if err != nil {
        responseError(w, "Cannot unmarshal response from login-api", err, http.StatusInternalServerError)
        return
    }

    if res.Error == "" {
        if res.Result == "Invalid username" || res.Result == "Invalid password" {
            responseError(w, "Invalid username or password", errors.New("Invalid username or password"), http.StatusBadRequest)
            return
        }

        // TODO: Add token in redis DB
        token := res.Result
        log.Println("Add token in DB")
        log.Println(token)

        responseJSON(w, "Login successful", "Login successful")
        return
    }

    responseError(w, "Cannot loggin user", nil, http.StatusInternalServerError)
    return
}

////////////////////////////////////////////////////
// Operations api
////////////////////////////////////////////////////
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: create subscription")

    // TODO: Check authentication token

    // Read body
    params := mux.Vars(r)
    id := params["username"]
    body, _ := ioutil.ReadAll(r.Body)

    //Sending request to operations-api
    response, err := http.Post("http://operations-api:5000/subscriptions/" + id, "application/json", bytes.NewReader(body))
    if err != nil {
        responseError(w, "Cannot send request to operations-api", err, http.StatusInternalServerError)
        return
    }

    // Unmarshal response
    var res model.ResponseResult
    body, _ = ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &res)
    if err != nil {
        responseError(w, "Cannot unmarshal response from operations-api", err, http.StatusInternalServerError)
        return
    }


    if res.Error == "" {
        responseJSON2(w, "Operation successful", "Operation successful", res.Data)
        return
    }

    responseError(w, "Operation failed", nil, http.StatusInternalServerError)
    return
}

func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: get all subscriptions")

    // TODO: Check authentication token

    // Read body
    params := mux.Vars(r)
    id := params["username"]
    body, _ := ioutil.ReadAll(r.Body)

    // Sending request to operations-api
    response, err := http.Get("http://operations-api:5000/subscriptions/" + id)
    if err != nil {
        responseError(w, "Cannot send request to operations-api", err, http.StatusInternalServerError)
        return
    }

    // Unmarshal response
    var res model.ResponseResult
    body, _ = ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &res)
    if err != nil {
        responseError(w, "Cannot unmarshal response from operations-api", err, http.StatusInternalServerError)
        return
    }

    if res.Error == "" {
        responseJSON2(w, "Operation successful", "Operation successful", res.Data)
        return
    }

    responseError(w, "Operation failed", nil, http.StatusInternalServerError)
    return
}

func GetSubscription(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: get one subscription")

    // TODO: Check authentication token

    // Read body
    params := mux.Vars(r)
    id := params["username"]
    subscriptionName := params["subName"]
    body, _ := ioutil.ReadAll(r.Body)

    // Sending request to operations-api
    response, err := http.Get("http://operations-api:5000/subscription/" + id + "/" + subscriptionName)
    if err != nil {
        responseError(w, "Cannot unmarshal response from operations-api", err, http.StatusInternalServerError)
        return
    }

    // Unmarshal response
    var res model.ResponseResult
    body, _ = ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &res)
    if err != nil {
        responseError(w, "Cannot unmarshal response from operations-api", err, http.StatusInternalServerError)
        return
    }

    if res.Error == "" {
        responseJSON2(w, "Operation successful", "Operation successful", res.Data)
        return
    }

    responseError(w, "Operation failed", nil, http.StatusInternalServerError)
    return
}


////////////////////////////////////////////////////
// Util functions
////////////////////////////////////////////////////
func responseError(w http.ResponseWriter, logMessage string, err error, code int) {
    log.Println("Error: " + logMessage)
    if err != nil {
        log.Println(err.Error())
    }
    log.Println("here")

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    var res model.ResponseResult
    // res.Error = err.Error()
    res.Error = logMessage
    json.NewEncoder(w).Encode(res)
}

func responseJSON(w http.ResponseWriter, logMessage string, result string) {
    log.Println(logMessage)
    w.Header().Set("Content-Type", "application/json")
    var res model.ResponseResult
    res.Result = result
    json.NewEncoder(w).Encode(res)
}

func responseJSON2(w http.ResponseWriter, logMessage string, result string, data []model.Subscription) {
    log.Println(logMessage)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    var res model.ResponseResult
    res.Result = result
    res.Data = data
    json.NewEncoder(w).Encode(res)
}