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

    // jwt "github.com/dgrijalva/jwt-go"
    // "go.mongodb.org/mongo-driver/bson"
    // "golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request")

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

    var res model.ResponseResult
    body, _ = ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &res)
    if err != nil {
        responseError(w, "Cannot unmarshal response from login-api", err, http.StatusInternalServerError)
        return
    }

    if res.Result == "Registration successful" {
        // Sending request to operations-api
        user.Password = ""
        user.Token = ""
        jsonValue, _ := json.Marshal(user)
        response, err := http.Post("http://operations-api:5000/user", "application/json", bytes.NewBuffer(jsonValue))
        if err != nil {
            responseError(w, "Cannot send request to operations-api", err, http.StatusInternalServerError)
            return
        }

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
    log.Println("Received request")
    // w.Header().Set("Content-Type", "application/json")

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
        log.Println("Error in sending request to auth")
        log.Println(err)
        return
    }

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



func responseError(w http.ResponseWriter, logMessage string, err error, code int) {
    log.Println("Error: " + logMessage)
    log.Println(err.Error())

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    var res model.ResponseResult
    res.Error = err.Error()
    json.NewEncoder(w).Encode(res)
}

func responseJSON(w http.ResponseWriter, logMessage string, result string) {
    log.Println(logMessage)
    w.Header().Set("Content-Type", "application/json")
    var res model.ResponseResult
    res.Result = result
    json.NewEncoder(w).Encode(res)
}