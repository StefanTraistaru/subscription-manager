package controller

import (
    "encoding/json"
    "log"
    "fmt"
    "errors"
    "strings"
    "bytes"
    "../model"
    "io/ioutil"
    "net/http"

    "github.com/gorilla/mux"
    jwt "github.com/dgrijalva/jwt-go"
    // "go.mongodb.org/mongo-driver/bson"
    // "golang.org/x/crypto/bcrypt"
)

////////////////////////////////////////////////////
// Login api handlers
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

        // Create and send response
        if res.Result == "User created successfully" {
            var response model.ResponseResult
            response.Result = "Registration successful"
            responseJSON(w, "Registration successful", response)
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

        // Add token in redis DB
        token := res.Result
        log.Println("Add token in DB")
        log.Println(token)

        // Create and send response
        var response model.ResponseResult
        response.Result = token
        responseJSON(w, "Login successful", response)
        return
    }

    responseError(w, "Cannot loggin user", nil, http.StatusInternalServerError)
    return
}

func Delete(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: delete")

    // Read body
    params := mux.Vars(r)
    username := params["username"]
    // Read body of the request

    // Sending request to login-api
    response, err := http.Post("http://login-api:5002/delete/" + username, "application/json", nil)
    if err != nil {
        log.Println("Error in sending request to login-api")
        log.Println(err)
        return
    }

    // Unmarshall response from login-api
    var res model.ResponseResult
    body, _ := ioutil.ReadAll(response.Body)
    err = json.Unmarshal(body, &res)
    if err != nil {
        responseError(w, "Cannot unmarshal response from login-api", err, http.StatusInternalServerError)
        return
    }

    if res.Error == "" {
        if res.Result == "Invalid user" {
            responseError(w, "Invalid username", errors.New("Invalid username"), http.StatusBadRequest)
            return
        }

        // Delete user from data db
        response, err := http.Post("http://operations-api:5000/user/delete/" + username, "application/json", nil)
        if err != nil {
            log.Println("Error in sending request to operations-api")
            log.Println(err)
            return
        }

        // Unmarshall response from operations-api
        var res model.ResponseResult
        body, _ := ioutil.ReadAll(response.Body)
        err = json.Unmarshal(body, &res)
        if err != nil {
            responseError(w, "Cannot unmarshal response from operations-api", err, http.StatusInternalServerError)
            return
        }

        // Create and send response
        if res.Result == "Deletion successful" {
            var response model.ResponseResult
            response.Result = "Deletion successful"
            responseJSON(w, "Deletion successful", response)
            return
        }

        responseError(w, "Cannot delete user", nil, http.StatusInternalServerError)
        return
    }

    responseError(w, "Cannot delete user", nil, http.StatusInternalServerError)
    return
}

////////////////////////////////////////////////////
// Operations api handlers
////////////////////////////////////////////////////
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: create subscription")

    // Check authentication token
    tokenString := r.Header.Get("Authorization")
    log.Println(tokenString)
    tokenString = strings.Split(tokenString, "Bearer ")[1]
    user, ok := verifyToken(tokenString)
    if !ok {
        responseError(w, "Token is not valid", nil, http.StatusInternalServerError)
        return
    }
    log.Println(user)

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

    // Create and send response
    if res.Error == "" {
        var response model.ResponseResult
        response.Result = "Operation successful"
        response.Data = res.Data
        responseJSON(w, "Operation successful", response)
        return
    }

    responseError(w, "Operation failed", nil, http.StatusInternalServerError)
    return
}

func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: get all subscriptions")

    // Check authentication token
    tokenString := r.Header.Get("Authorization")
    log.Println(tokenString)
    tokenString = strings.Split(tokenString, "Bearer ")[1]
    user, ok := verifyToken(tokenString)
    if !ok {
        responseError(w, "Token is not valid", nil, http.StatusInternalServerError)
        return
    }
    log.Println(user)

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

    // Create and send response
    if res.Error == "" {
        var response model.ResponseResult
        response.Result = "Operation successful"
        response.Data = res.Data
        responseJSON(w, "Operation successful", response)
        return
    }

    responseError(w, "Operation failed", nil, http.StatusInternalServerError)
    return
}

func GetSubscription(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request: get one subscription")

    // Check authentication token
    tokenString := r.Header.Get("Authorization")
    log.Println(tokenString)
    tokenString = strings.Split(tokenString, "Bearer ")[1]
    user, ok := verifyToken(tokenString)
    if !ok {
        responseError(w, "Token is not valid", nil, http.StatusInternalServerError)
        return
    }
    log.Println(user)

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

    // Create and send response
    if res.Error == "" {
        var response model.ResponseResult
        response.Result = "Operation successful"
        response.Data = res.Data
        responseJSON(w, "Operation successful", response)
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

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    var res model.ResponseResult
    // res.Error = err.Error()
    res.Error = logMessage
    json.NewEncoder(w).Encode(res)
}

func responseJSON(w http.ResponseWriter, message string, data interface{}) {
    log.Println(message)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(data)
}

func verifyToken(tokenString string) (model.User, bool){
    log.Println(tokenString)
    var result model.User

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Validate the algorithm
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method")
        }
        return []byte("secret"), nil
    })
    if err != nil {
        log.Println("Error in jwt.Parse function")
        return result, false
    }

    // var res model.ResponseResult
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        result.Username = claims["username"].(string)
        result.FirstName = claims["firstname"].(string)
        result.LastName = claims["lastname"].(string)
        return result, true
    } else {
        return result, false
    }
}